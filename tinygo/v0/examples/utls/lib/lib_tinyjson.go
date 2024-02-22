// Code generated by tinyjson for marshaling/unmarshaling. DO NOT EDIT.

package lib

import (
	tinyjson "github.com/CosmWasm/tinyjson"
	jlexer "github.com/CosmWasm/tinyjson/jlexer"
	jwriter "github.com/CosmWasm/tinyjson/jwriter"
)

// suppress unused package warning
var (
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ tinyjson.Marshaler
)

func tinyjsonAded76e7DecodeGithubComRefractionNetworkingWatmTinygoV0ExamplesUtlsLib(in *jlexer.Lexer, out *TLSConfig) {
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
		case "next_protos":
			if in.IsNull() {
				in.Skip()
				out.NextProtos = nil
			} else {
				in.Delim('[')
				if out.NextProtos == nil {
					if !in.IsDelim(']') {
						out.NextProtos = make([]string, 0, 4)
					} else {
						out.NextProtos = []string{}
					}
				} else {
					out.NextProtos = (out.NextProtos)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.NextProtos = append(out.NextProtos, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "application_settings":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.ApplicationSettings = make(map[string][]uint8)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v2 []uint8
					if in.IsNull() {
						in.Skip()
						v2 = nil
					} else {
						v2 = in.Bytes()
					}
					(out.ApplicationSettings)[key] = v2
					in.WantComma()
				}
				in.Delim('}')
			}
		case "server_name":
			out.ServerName = string(in.String())
		case "insecure_skip_verify":
			out.InsecureSkipVerify = bool(in.Bool())
		case "insecure_skip_time_verify":
			out.InsecureSkipTimeVerify = bool(in.Bool())
		case "omit_empty_psk":
			out.OmitEmptyPsk = bool(in.Bool())
		case "insecure_server_name_to_verify":
			out.InsecureServerNameToVerify = string(in.String())
		case "session_tickets_disabled":
			out.SessionTicketsDisabled = bool(in.Bool())
		case "pq_signature_schemes_enabled":
			out.PQSignatureSchemesEnabled = bool(in.Bool())
		case "dynamic_record_sizing_disabled":
			out.DynamicRecordSizingDisabled = bool(in.Bool())
		case "ech_configs":
			if in.IsNull() {
				in.Skip()
				out.ECHConfigs = nil
			} else {
				out.ECHConfigs = in.Bytes()
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
func tinyjsonAded76e7EncodeGithubComRefractionNetworkingWatmTinygoV0ExamplesUtlsLib(out *jwriter.Writer, in TLSConfig) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"next_protos\":"
		out.RawString(prefix[1:])
		if in.NextProtos == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.NextProtos {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"application_settings\":"
		out.RawString(prefix)
		if in.ApplicationSettings == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v7First := true
			for v7Name, v7Value := range in.ApplicationSettings {
				if v7First {
					v7First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v7Name))
				out.RawByte(':')
				out.Base64Bytes(v7Value)
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"server_name\":"
		out.RawString(prefix)
		out.String(string(in.ServerName))
	}
	{
		const prefix string = ",\"insecure_skip_verify\":"
		out.RawString(prefix)
		out.Bool(bool(in.InsecureSkipVerify))
	}
	{
		const prefix string = ",\"insecure_skip_time_verify\":"
		out.RawString(prefix)
		out.Bool(bool(in.InsecureSkipTimeVerify))
	}
	{
		const prefix string = ",\"omit_empty_psk\":"
		out.RawString(prefix)
		out.Bool(bool(in.OmitEmptyPsk))
	}
	{
		const prefix string = ",\"insecure_server_name_to_verify\":"
		out.RawString(prefix)
		out.String(string(in.InsecureServerNameToVerify))
	}
	{
		const prefix string = ",\"session_tickets_disabled\":"
		out.RawString(prefix)
		out.Bool(bool(in.SessionTicketsDisabled))
	}
	{
		const prefix string = ",\"pq_signature_schemes_enabled\":"
		out.RawString(prefix)
		out.Bool(bool(in.PQSignatureSchemesEnabled))
	}
	{
		const prefix string = ",\"dynamic_record_sizing_disabled\":"
		out.RawString(prefix)
		out.Bool(bool(in.DynamicRecordSizingDisabled))
	}
	{
		const prefix string = ",\"ech_configs\":"
		out.RawString(prefix)
		out.Base64Bytes(in.ECHConfigs)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v TLSConfig) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjsonAded76e7EncodeGithubComRefractionNetworkingWatmTinygoV0ExamplesUtlsLib(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v TLSConfig) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjsonAded76e7EncodeGithubComRefractionNetworkingWatmTinygoV0ExamplesUtlsLib(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TLSConfig) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjsonAded76e7DecodeGithubComRefractionNetworkingWatmTinygoV0ExamplesUtlsLib(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *TLSConfig) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjsonAded76e7DecodeGithubComRefractionNetworkingWatmTinygoV0ExamplesUtlsLib(l, v)
}
func tinyjsonAded76e7DecodeGithubComRefractionNetworkingWatmTinygoV0ExamplesUtlsLib1(in *jlexer.Lexer, out *Configurables) {
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
		case "tls_config":
			if in.IsNull() {
				in.Skip()
				out.TLSConfig = nil
			} else {
				if out.TLSConfig == nil {
					out.TLSConfig = new(TLSConfig)
				}
				(*out.TLSConfig).UnmarshalTinyJSON(in)
			}
		case "client_hello_id":
			out.ClientHelloID = string(in.String())
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
func tinyjsonAded76e7EncodeGithubComRefractionNetworkingWatmTinygoV0ExamplesUtlsLib1(out *jwriter.Writer, in Configurables) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"tls_config\":"
		out.RawString(prefix[1:])
		if in.TLSConfig == nil {
			out.RawString("null")
		} else {
			(*in.TLSConfig).MarshalTinyJSON(out)
		}
	}
	{
		const prefix string = ",\"client_hello_id\":"
		out.RawString(prefix)
		out.String(string(in.ClientHelloID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Configurables) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	tinyjsonAded76e7EncodeGithubComRefractionNetworkingWatmTinygoV0ExamplesUtlsLib1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalTinyJSON supports tinyjson.Marshaler interface
func (v Configurables) MarshalTinyJSON(w *jwriter.Writer) {
	tinyjsonAded76e7EncodeGithubComRefractionNetworkingWatmTinygoV0ExamplesUtlsLib1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Configurables) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	tinyjsonAded76e7DecodeGithubComRefractionNetworkingWatmTinygoV0ExamplesUtlsLib1(&r, v)
	return r.Error()
}

// UnmarshalTinyJSON supports tinyjson.Unmarshaler interface
func (v *Configurables) UnmarshalTinyJSON(l *jlexer.Lexer) {
	tinyjsonAded76e7DecodeGithubComRefractionNetworkingWatmTinygoV0ExamplesUtlsLib1(l, v)
}
