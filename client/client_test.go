package client

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetRemoteFileWithoutDecryption(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash := NewMockHash(ctrl)
	store := NewMockStore(ctrl)

	c := &Client{
		Hash:  hash,
		Store: store,
	}

	store.EXPECT().DownloadFile("s3://foo/bar.txt", "bar.txt").Return("muahahaha", nil)
	hash.EXPECT().Verify("bar.txt", "muahahaha").Return(nil)

	assert.NoError(t, c.GetRemoteFile("s3://foo/bar.txt", "bar.txt"))
}

func TestGetRemoteFileWithDecryption(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash := NewMockHash(ctrl)
	store := NewMockStore(ctrl)
	cipher := NewMockCipher(ctrl)

	c := &Client{
		Hash:   hash,
		Store:  store,
		Cipher: cipher,
	}

	store.EXPECT().DownloadFile("s3://foo/bar.txt", "bar.txt.tmp").Return("muahahaha", nil)
	hash.EXPECT().Verify("bar.txt.tmp", "muahahaha").Return(nil)
	cipher.EXPECT().Decrypt("bar.txt.tmp", "bar.txt").Return(nil)

	assert.NoError(t, c.GetRemoteFile("s3://foo/bar.txt", "bar.txt"))
}

func TestPutLocalFileWithoutEncryption(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash := NewMockHash(ctrl)
	store := NewMockStore(ctrl)

	c := &Client{
		Hash:  hash,
		Store: store,
	}

	hash.EXPECT().Calculate("bar.txt").Return("woahahaha", nil)
	store.EXPECT().UploadFile("s3://foo/bar.txt", "bar.txt", "woahahaha").Return(nil)

	assert.NoError(t, c.PutLocalFile("s3://foo/bar.txt", "bar.txt"))
}

func TestPutLocalFileWithEncryption(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	hash := NewMockHash(ctrl)
	store := NewMockStore(ctrl)
	cipher := NewMockCipher(ctrl)

	c := &Client{
		Hash:   hash,
		Store:  store,
		Cipher: cipher,
	}

	cipher.EXPECT().Encrypt("bar.txt", "bar.txt.tmp").Return(nil)
	hash.EXPECT().Calculate("bar.txt.tmp").Return("woahahaha", nil)
	store.EXPECT().UploadFile("s3://foo/bar.txt", "bar.txt.tmp", "woahahaha").Return(nil)

	assert.NoError(t, c.PutLocalFile("s3://foo/bar.txt", "bar.txt"))
}
