// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package rest

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson84463b63DecodeGoRestInternalRest(in *jlexer.Lexer, out *UserForm) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "email":
			out.Email = string(in.String())
		case "pass":
			out.Pass = string(in.String())
		case "repeat_pass":
			out.RepeatPass = string(in.String())
		case "is_admin":
			out.IsAdmin = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson84463b63EncodeGoRestInternalRest(out *jwriter.Writer, in UserForm) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"email\":"
		out.RawString(prefix[1:])
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"pass\":"
		out.RawString(prefix)
		out.String(string(in.Pass))
	}
	{
		const prefix string = ",\"repeat_pass\":"
		out.RawString(prefix)
		out.String(string(in.RepeatPass))
	}
	{
		const prefix string = ",\"is_admin\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsAdmin))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserForm) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson84463b63EncodeGoRestInternalRest(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserForm) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson84463b63EncodeGoRestInternalRest(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserForm) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson84463b63DecodeGoRestInternalRest(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserForm) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson84463b63DecodeGoRestInternalRest(l, v)
}
func easyjson84463b63DecodeGoRestInternalRest1(in *jlexer.Lexer, out *TaskForm) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "text":
			out.Text = string(in.String())
		case "date":
			out.Date = string(in.String())
		case "tags":
			if in.IsNull() {
				in.Skip()
				out.Tags = nil
			} else {
				in.Delim('[')
				if out.Tags == nil {
					if !in.IsDelim(']') {
						out.Tags = make([]*TagForm, 0, 8)
					} else {
						out.Tags = []*TagForm{}
					}
				} else {
					out.Tags = (out.Tags)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *TagForm
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(TagForm)
						}
						(*v1).UnmarshalEasyJSON(in)
					}
					out.Tags = append(out.Tags, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "done":
			out.Done = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson84463b63EncodeGoRestInternalRest1(out *jwriter.Writer, in TaskForm) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix[1:])
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.String(string(in.Date))
	}
	{
		const prefix string = ",\"tags\":"
		out.RawString(prefix)
		if in.Tags == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Tags {
				if v2 > 0 {
					out.RawByte(',')
				}
				if v3 == nil {
					out.RawString("null")
				} else {
					(*v3).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"done\":"
		out.RawString(prefix)
		out.Bool(bool(in.Done))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v TaskForm) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson84463b63EncodeGoRestInternalRest1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TaskForm) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson84463b63EncodeGoRestInternalRest1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TaskForm) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson84463b63DecodeGoRestInternalRest1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TaskForm) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson84463b63DecodeGoRestInternalRest1(l, v)
}
func easyjson84463b63DecodeGoRestInternalRest2(in *jlexer.Lexer, out *TagForm) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "text":
			out.Text = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson84463b63EncodeGoRestInternalRest2(out *jwriter.Writer, in TagForm) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix[1:])
		out.String(string(in.Text))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v TagForm) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson84463b63EncodeGoRestInternalRest2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TagForm) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson84463b63EncodeGoRestInternalRest2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TagForm) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson84463b63DecodeGoRestInternalRest2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TagForm) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson84463b63DecodeGoRestInternalRest2(l, v)
}
