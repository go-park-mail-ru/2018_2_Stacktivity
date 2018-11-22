// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

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

func easyjson85f0d656Decode20182StacktivityModels(in *jlexer.Lexer, out *Message) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "event":
			out.Event = int(in.Int())
		case "players":
			if in.IsNull() {
				in.Skip()
				out.Players = nil
			} else {
				if out.Players == nil {
					out.Players = new([]string)
				}
				if in.IsNull() {
					in.Skip()
					*out.Players = nil
				} else {
					in.Delim('[')
					if *out.Players == nil {
						if !in.IsDelim(']') {
							*out.Players = make([]string, 0, 4)
						} else {
							*out.Players = []string{}
						}
					} else {
						*out.Players = (*out.Players)[:0]
					}
					for !in.IsDelim(']') {
						var v1 string
						v1 = string(in.String())
						*out.Players = append(*out.Players, v1)
						in.WantComma()
					}
					in.Delim(']')
				}
			}
		case "level":
			if in.IsNull() {
				in.Skip()
				out.Level = nil
			} else {
				if out.Level == nil {
					out.Level = new(Level)
				}
				(*out.Level).UnmarshalEasyJSON(in)
			}
		case "curve":
			if in.IsNull() {
				in.Skip()
				out.Curve = nil
			} else {
				if out.Curve == nil {
					out.Curve = new([]Dot)
				}
				if in.IsNull() {
					in.Skip()
					*out.Curve = nil
				} else {
					in.Delim('[')
					if *out.Curve == nil {
						if !in.IsDelim(']') {
							*out.Curve = make([]Dot, 0, 4)
						} else {
							*out.Curve = []Dot{}
						}
					} else {
						*out.Curve = (*out.Curve)[:0]
					}
					for !in.IsDelim(']') {
						var v2 Dot
						(v2).UnmarshalEasyJSON(in)
						*out.Curve = append(*out.Curve, v2)
						in.WantComma()
					}
					in.Delim(']')
				}
			}
		case "ball":
			if in.IsNull() {
				in.Skip()
				out.Ball = nil
			} else {
				if out.Ball == nil {
					out.Ball = new(Ball)
				}
				(*out.Ball).UnmarshalEasyJSON(in)
			}
		case "status":
			if in.IsNull() {
				in.Skip()
				out.Status = nil
			} else {
				if out.Status == nil {
					out.Status = new(string)
				}
				*out.Status = string(in.String())
			}
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
func easyjson85f0d656Encode20182StacktivityModels(out *jwriter.Writer, in Message) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"event\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Event))
	}
	if in.Players != nil {
		const prefix string = ",\"players\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if *in.Players == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v3, v4 := range *in.Players {
				if v3 > 0 {
					out.RawByte(',')
				}
				out.String(string(v4))
			}
			out.RawByte(']')
		}
	}
	if in.Level != nil {
		const prefix string = ",\"level\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Level).MarshalEasyJSON(out)
	}
	if in.Curve != nil {
		const prefix string = ",\"curve\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if *in.Curve == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range *in.Curve {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	if in.Ball != nil {
		const prefix string = ",\"ball\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Ball).MarshalEasyJSON(out)
	}
	if in.Status != nil {
		const prefix string = ",\"status\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(*in.Status))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Message) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656Encode20182StacktivityModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Message) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656Encode20182StacktivityModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Message) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656Decode20182StacktivityModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Message) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656Decode20182StacktivityModels(l, v)
}
func easyjson85f0d656Decode20182StacktivityModels1(in *jlexer.Lexer, out *Level) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "levelNumber":
			out.LevelNumber = int(in.Int())
		case "balls":
			if in.IsNull() {
				in.Skip()
				out.Balls = nil
			} else {
				in.Delim('[')
				if out.Balls == nil {
					if !in.IsDelim(']') {
						out.Balls = make([]Ball, 0, 1)
					} else {
						out.Balls = []Ball{}
					}
				} else {
					out.Balls = (out.Balls)[:0]
				}
				for !in.IsDelim(']') {
					var v7 Ball
					(v7).UnmarshalEasyJSON(in)
					out.Balls = append(out.Balls, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson85f0d656Encode20182StacktivityModels1(out *jwriter.Writer, in Level) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"levelNumber\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.LevelNumber))
	}
	{
		const prefix string = ",\"balls\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		if in.Balls == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Balls {
				if v8 > 0 {
					out.RawByte(',')
				}
				(v9).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Level) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656Encode20182StacktivityModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Level) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656Encode20182StacktivityModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Level) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656Decode20182StacktivityModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Level) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656Decode20182StacktivityModels1(l, v)
}
func easyjson85f0d656Decode20182StacktivityModels2(in *jlexer.Lexer, out *Dot) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "x":
			out.X = int(in.Int())
		case "y":
			out.Y = int(in.Int())
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
func easyjson85f0d656Encode20182StacktivityModels2(out *jwriter.Writer, in Dot) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"x\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.X))
	}
	{
		const prefix string = ",\"y\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Y))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Dot) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656Encode20182StacktivityModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Dot) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656Encode20182StacktivityModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Dot) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656Decode20182StacktivityModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Dot) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656Decode20182StacktivityModels2(l, v)
}
func easyjson85f0d656Decode20182StacktivityModels3(in *jlexer.Lexer, out *Ball) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "number":
			out.Number = int(in.Int())
		case "x":
			out.X = int(in.Int())
		case "y":
			out.Y = int(in.Int())
		case "r":
			out.R = int(in.Int())
		case "type":
			out.Type = string(in.String())
		case "color":
			out.Color = string(in.String())
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
func easyjson85f0d656Encode20182StacktivityModels3(out *jwriter.Writer, in Ball) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"number\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Number))
	}
	{
		const prefix string = ",\"x\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.X))
	}
	{
		const prefix string = ",\"y\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.Y))
	}
	{
		const prefix string = ",\"r\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.R))
	}
	{
		const prefix string = ",\"type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"color\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Color))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Ball) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson85f0d656Encode20182StacktivityModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Ball) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson85f0d656Encode20182StacktivityModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Ball) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson85f0d656Decode20182StacktivityModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Ball) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson85f0d656Decode20182StacktivityModels3(l, v)
}
