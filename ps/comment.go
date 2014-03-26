package ps

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

type comment struct {
	raw             bool
	key, value, out string
}

func (c comment) String() string {
	if c.raw {
		return c.value
	}
	if c.out != "" {
		return c.out
	}
	if c.value != "" {
		c.out = fmt.Sprintf("%%%%%s: %s", c.key, c.value)
	} else {
		c.out = "%%" + c.key
	}
	return c.out
}

var forbidden = []string{
	//Header section
	"Creator",
	"CreationDate",
	"Pages",
	"BoundingBox",
	"DocumentData",
	"LanguageLevel",
	"EndComments",
	//Setup section
	"BeginSetup",
	"EndSetup",
	//Page setup section
	"BeginPageSetup",
	"PageBoundingBox",
	"EndPageSetup",
	//Other sections
	"BeginProlog",
	"EndProlog",
	"Page",
	"Trailer",
	"EOF",
}

func (c comment) invalidKey() error {
	if c.raw {
		return nil
	}
	if c.key == "" {
		return errors.New("No header specified")
	}
	if strings.IndexFunc(c.key, unicode.IsSpace) != -1 {
		return fmt.Errorf("DSC key cannot contain spaces, got: %s", c.key)
	}
	if strings.IndexRune(c.key, ':') != -1 {
		return fmt.Errorf("DSC key cannot contain colon, got: %s", c.key)
	}
	for _, k := range forbidden {
		if k == c.key {
			return fmt.Errorf("use of comment type %s is forbidden", k)
		}
	}
	return nil
}

func (c comment) invalidValue() error {
	if c.raw {
		return nil
	}
	if c.value != "" && strings.IndexAny(c.value, "\n\r") != -1 {
		return errors.New("DSC key cannot contain newlines")
	}
	return nil
}

func (c comment) invalidLength() error {
	if c.raw {
		return nil
	}
	s := c.String()
	if ln := len([]byte(s)); ln > 255 {
		return fmt.Errorf("comment exceeds maximum length of 255, got %d", ln)
	}
	return nil
}

func (c comment) Err() (err error) {
	if c.raw {
		return nil
	}
	if err = c.invalidKey(); err != nil {
		return err
	}
	if err = c.invalidValue(); err != nil {
		return err
	}
	return c.invalidLength()
}

//Comments represents a sequence of PostScript Document Structuring Comments
//(DSC).
//
//Please see that manual for details on the available comments and their
//meanings.
//
//In particular, the %%IncludeFeature comment allows a device-independent means
//of controlling printer device features, so the PostScript Printer Description
//Files Specification will also be a useful reference.
//
//The individual comment entries have String and Err methods.
//
//See Comment for an explanation for how a comment is constructed and what
//constitutes a valid comment.
//
//Comments can be reused and do not need to be reconstituted on each use
//if they do not change.
type Comments []comment

//Err calls Err on each comment in turn and returns the first error found.
func (c Comments) Err() error {
	for _, c := range c {
		if err := c.Err(); err != nil {
			return err
		}
	}
	return nil
}

//Comment specifies a PostScript Document Structuring Comment (DSC).
//
//The returned comment has String and Err methods.
//
//The returned comment's String method produces this output:
//	%%key: value
//
//If value == "", String produces:
//	%%key
//
//The total byte length of the result of String must not be greater than 255.
//
//As a convenience, Comment trims trailing and leading whitespace from the key
//and value; however, it is an error for a key to contain any other whitespace
//and for value to contain any newlines after being trimmed.
//
//It is also an error for key to contain the ":" colon character.
//
//The following keys are used by libcairo and are forbidden:
//	Creator
//	CreationDate
//	Pages
//	BoundingBox
//	DocumentData
//	LanguageLevel
//	EndComments
//	BeginSetup
//	EndSetup
//	BeginProlog
//	EndProlog
//	Page
//	Trailer
//	EOF
//
//Use of a forbidden key results in an error.
//
//Even if a comment does not result in an error, that does not mean it produces
//the desires effect, only that it is not invalid.
func Comment(key, value string) comment {
	return comment{
		key:   strings.TrimSpace(key),
		value: strings.TrimSpace(value),
	}
}

//Commentf is a convenience function for
//	Comment(key, fmt.Sprintf(value, vars...))
//
//See Comment for an explanation for how a comment is constructed and what
//constitutes a valid comment.
func Commentf(key, value string, vars ...interface{}) comment {
	return Comment(key, fmt.Sprintf(value, vars...))
}

//RawComment creates a comments that is exempt from formatting and error
//checking.
//Use at your own peril.
func RawComment(s string) comment {
	return comment{
		raw:   true,
		value: s,
	}
}
