package ssh

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"sync"
	"time"
)

type PassThru func(r io.Reader, total int64) io.Reader

func (a *ConnInfo) Upload(
	ctx context.Context,
	fileReader io.Reader,
	remotePath string,
	permissions string,
) error {
	return a.CopyFilePassThru(ctx, fileReader, remotePath, permissions, nil)
}

func (a *ConnInfo) CopyFromFile(
	ctx context.Context,
	file os.File,
	remotePath string,
	permissions string,
) error {
	return a.CopyFromFilePassThru(ctx, file, remotePath, permissions, nil)
}

func (a *ConnInfo) CopyFromFilePassThru(
	ctx context.Context,
	file os.File,
	remotePath string,
	permissions string,
	passThru PassThru,
) error {
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	return a.Scp(ctx, &file, remotePath, permissions, stat.Size(), passThru)
}

func (c *ConnInfo) CopyFilePassThru(
	ctx context.Context,
	fileReader io.Reader,
	remotePath string,
	permissions string,
	passThru PassThru,
) error {
	contentsBytes, err := io.ReadAll(fileReader)
	if err != nil {
		return fmt.Errorf("failed to read all data from reader: %w", err)
	}
	bytesReader := bytes.NewReader(contentsBytes)

	return c.Scp(
		ctx,
		bytesReader,
		remotePath,
		permissions,
		int64(len(contentsBytes)),
		passThru,
	)
}

func (c *ConnInfo) Scp(
	ctx context.Context,
	r io.Reader,
	remotePath string,
	permissions string,
	size int64,
	passThru PassThru,
) error {
	if c.Client == nil {
		if _, err := c.NewClient(); err != nil {
			return err
		}
	}

	session, err := c.Client.NewSession()
	if err != nil {
		return fmt.Errorf("error creating ssh session in copy to remote: %v", err)
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	w, err := session.StdinPipe()
	if err != nil {
		return err
	}

	if passThru != nil {
		r = passThru(r, size)
	}

	filename := path.Base(remotePath)
	err = session.Start(fmt.Sprintf("%s -qt %q", "scp", remotePath))
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	errCh := make(chan error, 2)

	go func() {
		defer wg.Done()

		_, err := fmt.Fprintln(w, "C"+permissions, size, filename)
		if err != nil {
			errCh <- err
			return
		}

		if err = checkResponse(stdout); err != nil {
			errCh <- err
			return
		}

		_, err = io.Copy(w, r)
		if err != nil {
			errCh <- err
			return
		}

		_, err = fmt.Fprint(w, "\x00")
		if err != nil {
			errCh <- err
			return
		}

		if err = checkResponse(stdout); err != nil {
			errCh <- err
			return
		}

		w.Close() // Close the write pipe here
	}()

	// Wait for the process to exit
	go func() {
		defer wg.Done()
		err := session.Wait()
		if err != nil {
			errCh <- err
			return
		}
	}()

	// If there is a timeout, stop the transfer if it has been exceeded
	if c.DialTimeOut == 0 {
		c.DialTimeOut = 5 * time.Second
	}
	if c.DialTimeOut > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.DialTimeOut)
		defer cancel()
	}

	// Wait for one of the conditions (error/timeout/completion) to occur
	if err := wait(&wg, ctx); err != nil {
		return err
	}

	close(errCh)

	// Collect any errors from the error channel
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func wait(wg *sync.WaitGroup, ctx context.Context) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func checkResponse(r io.Reader) error {
	_, err := ParseResponse(r, nil)
	if err != nil {
		return err
	}

	return nil

}

func (a *ConnInfo) CopyFromRemote(ctx context.Context, file *os.File, remotePath string) error {
	return a.CopyFromRemotePassThru(ctx, file, remotePath, nil)
}

// CopyFromRemotePassThru copies a file from the remote to the given writer. The passThru parameter can be used
// to keep track of progress and how many bytes that were download from the remote.
// `passThru` can be set to nil to disable this behaviour.
func (a *ConnInfo) CopyFromRemotePassThru(
	ctx context.Context,
	w io.Writer,
	remotePath string,
	passThru PassThru,
) error {
	_, err := a.copyFromRemote(ctx, w, remotePath, passThru, false)

	return err
}

// CopyFroRemoteFileInfos copies a file from the remote to a given writer and return a FileInfos struct
// containing information about the file such as permissions, the file size, modification time and access time
func (a *ConnInfo) CopyFromRemoteFileInfos(
	ctx context.Context,
	w io.Writer,
	remotePath string,
	passThru PassThru,
) (*FileInfos, error) {
	return a.copyFromRemote(ctx, w, remotePath, passThru, true)
}

func (a *ConnInfo) copyFromRemote(
	ctx context.Context,
	w io.Writer,
	remotePath string,
	passThru PassThru,
	preserveFileTimes bool,
) (*FileInfos, error) {
	if a.Client == nil {
		if _, err := a.NewClient(); err != nil {
			return nil, err
		}
	}

	session, err := a.Client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Error creating ssh session in copy from remote: %v", err)
	}
	defer a.Client.Close()

	wg := sync.WaitGroup{}
	errCh := make(chan error, 4)
	var fileInfos *FileInfos

	wg.Add(1)
	go func() {
		var err error

		defer func() {
			// NOTE: this might send an already sent error another time, but since we only receive one, this is fine. On the "happy-path" of this function, the error will be `nil` therefore completing the "err<-errCh" at the bottom of the function.
			errCh <- err
			// We must unblock the go routine first as we block on reading the channel later
			wg.Done()

		}()

		r, err := session.StdoutPipe()
		if err != nil {
			errCh <- err
			return
		}

		in, err := session.StdinPipe()
		if err != nil {
			errCh <- err
			return
		}
		defer in.Close()

		if preserveFileTimes {
			err = session.Start(fmt.Sprintf("%s -pf %q", "scp", remotePath))
		} else {
			err = session.Start(fmt.Sprintf("%s -f %q", "scp", remotePath))
		}
		if err != nil {
			errCh <- err
			return
		}

		err = Ack(in)
		if err != nil {
			errCh <- err
			return
		}

		fileInfo, err := ParseResponse(r, in)
		if err != nil {
			errCh <- err
			return
		}

		fileInfos = fileInfo

		err = Ack(in)
		if err != nil {
			errCh <- err
			return
		}

		if passThru != nil {
			r = passThru(r, fileInfo.Size)
		}

		_, err = CopyN(w, r, fileInfo.Size)
		if err != nil {
			errCh <- err
			return
		}

		err = Ack(in)
		if err != nil {
			errCh <- err
			return
		}

		err = session.Wait()
		if err != nil {
			errCh <- err
			return
		}
	}()

	if a.DialTimeOut == 0 {
		a.DialTimeOut = 5 * time.Second
	}

	if a.DialTimeOut > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, a.DialTimeOut)
		defer cancel()
	}

	if err := wait(&wg, ctx); err != nil {
		return nil, err
	}

	finalErr := <-errCh
	close(errCh)
	return fileInfos, finalErr
}

func CopyN(writer io.Writer, src io.Reader, size int64) (int64, error) {
	var total int64
	total = 0
	for total < size {
		n, err := io.CopyN(writer, src, size)
		if err != nil {
			return 0, err
		}
		total += n
	}

	return total, nil
}
