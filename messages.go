package main

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ParsedLine struct {
	MsgType          string
	TransType        sql.NullString
	SessionID        sql.NullString
	AircraftID       string
	HexIdent         sql.NullString
	FlightID         string
	DateMsgGenerated time.Time
	TimeMsgGenerated time.Time
	DateMsgLogged    time.Time
	TimeMsgLogged    time.Time
	Callsign         sql.NullString
	Altitude         sql.NullInt32
	GroundSpeed      sql.NullFloat64
	Track            sql.NullFloat64
	Longitude        sql.NullFloat64
	Latitude         sql.NullFloat64
	VerticalRate     sql.NullFloat64
	Squawk           sql.NullString
	Alert            sql.NullBool // indicates squawk has changed
	Emergency        sql.NullBool
	Ident            sql.NullBool
	IsOnGround       sql.NullBool
}

func Parse(rawLine string) (ParsedLine, error) {
	parts := strings.Split(rawLine, ",")
	if len(parts) < 2 {
		return ParsedLine{}, errors.New("less than two elements in raw Line")
	}
	colCount := len(parts) - 1
	var pl ParsedLine

	pl.MsgType = parts[0]
	pl.TransType = stringToSqlNullString(parts[1])
	pl.SessionID = stringToSqlNullString(parts[2])
	pl.AircraftID = parts[3]
	pl.HexIdent = stringToSqlNullString(parts[4])
	pl.FlightID = parts[5]
	pl.DateMsgGenerated = formatDate(parts[6])
	pl.TimeMsgGenerated = formatTime(parts[7])
	pl.DateMsgLogged = formatDate(parts[8])
	pl.TimeMsgLogged = formatTime(parts[9])
	if colCount < 10 {
		return pl, nil
	}
	pl.Callsign = stringToSqlNullString(parts[10])
	if colCount < 11 {
		return pl, nil
	}
	pl.Altitude = stringToSqlNullInt32(parts[11])
	pl.GroundSpeed = stringToSqlNullFloat64(parts[12])
	pl.Track = stringToSqlNullFloat64(parts[13])
	pl.Longitude = stringToSqlNullFloat64(parts[14])
	pl.Latitude = stringToSqlNullFloat64(parts[15])
	pl.VerticalRate = stringToSqlNullFloat64(parts[16])
	pl.Squawk = stringToSqlNullString(parts[17])
	pl.Alert = stringToSqlNullBool(parts[18])
	pl.Emergency = stringToSqlNullBool(parts[19])
	pl.Ident = stringToSqlNullBool(parts[20])
	pl.IsOnGround = stringToSqlNullBool(parts[21])
	return pl, nil

	return ParsedLine{}, nil
}

func stringToSqlNullString(s string) sql.NullString {
	var ns sql.NullString
	if s != "" {
		ns.String = s
		ns.Valid = true
	}
	return ns
}

func stringToSqlNullInt32(s string) sql.NullInt32 {
	var ns sql.NullInt32
	if s != "" {
		i, err := strconv.Atoi(s)
		if err == nil {
			ns.Int32 = int32(i)
			ns.Valid = true
		}
	}
	return ns
}

func stringToSqlNullFloat64(s string) sql.NullFloat64 {
	var ns sql.NullFloat64
	if s != "" {
		f, err := strconv.ParseFloat(s, 64)
		if err == nil {
			ns.Float64 = f
			ns.Valid = true
		}
	}
	return ns
}

func stringToSqlNullBool(s string) sql.NullBool {
	var ns sql.NullBool
	switch strings.ToLower(s) {
	case "0":
		ns.Bool = true
		ns.Valid = true
	default:
		ns.Valid = false
	}
	return ns
}
func formatDate(ts string) time.Time {
	tm, err := time.Parse("2006/01/02", ts)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return time.Time{}
	}
	fmt.Println(tm.String())
	return tm
}

func formatTime(ts string) time.Time {
	tm, err := time.Parse("15:04:05", ts)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return time.Time{}
	}
	fmt.Println("time: ", tm.String())

	return tm
}
