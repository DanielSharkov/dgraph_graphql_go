package apitest

import (
	"testing"

	"github.com/romshark/dgraph_graphql_go/api/graph/gqlmod"
	"github.com/romshark/dgraph_graphql_go/apitest/setup"
	"github.com/romshark/dgraph_graphql_go/store/errors"
)

// TestEditPostAuth tests post editing authorization
func TestEditPostAuth(t *testing.T) {
	setupTest := func(t *testing.T) (
		ts *setup.TestSetup,
		author *gqlmod.User,
		authorClt *setup.Client,
		post *gqlmod.Post,
	) {
		ts = setup.New(t, tcx)
		debug := ts.Debug()

		authorEmail := "author@tst.tst"
		authorPass := "testpass"
		author = debug.Help.OK.CreateUser(
			"fooBarowich",
			authorEmail,
			authorPass,
		)
		authorClt, _ = ts.Client(authorEmail, authorPass)

		post = debug.Help.OK.CreatePost(
			*author.ID,
			"example title",
			"example contents",
		)

		return
	}

	// Test editing posts as a guest
	t.Run("guest (noauth)", func(t *testing.T) {
		ts, author, _, post := setupTest(t)
		defer ts.Teardown()

		newTitle := "new test post"
		newContents := "new test content"
		ts.Guest().Help.ERR.EditPost(
			errors.ErrUnauthorized,
			*post.ID,
			*author.ID,
			&newTitle,
			&newContents,
		)
	})

	// Test editing posts on behalf of other users
	t.Run("non-editor (noauth)", func(t *testing.T) {
		ts, _, authorClt, post := setupTest(t)
		defer ts.Teardown()

		other := ts.Debug().Help.OK.CreateUser("other", "2@tst.tst", "testpass")

		newTitle := "new test post"
		newContents := "new test content"
		authorClt.Help.ERR.EditPost(
			errors.ErrUnauthorized,
			*post.ID,
			*other.ID, // Different editor ID
			&newTitle,
			&newContents,
		)
	})

	// Test editing posts of other users
	t.Run("non-author (noauth)", func(t *testing.T) {
		ts, _, _, post := setupTest(t)
		defer ts.Teardown()

		other := ts.Debug().Help.OK.CreateUser("other", "2@tst.tst", "testpass")
		otherClt, _ := ts.Client("2@tst.tst", "testpass")

		newTitle := "new test post"
		newContents := "new test content"
		otherClt.Help.ERR.EditPost(
			errors.ErrUnauthorized,
			*post.ID, // Someone else's post
			*other.ID,
			&newTitle,
			&newContents,
		)
	})
}
