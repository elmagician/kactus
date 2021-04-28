package types

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
)

// SetTyped matches string value to interface receiver type.
func SetTyped(value string, to interface{}) error { // nolint: gocyclo
	baseMsg := "Setting value to"

	switch typedTo := to.(type) {
	case *string:
		log.Debug(baseMsg + " string")

		*typedTo = value
	case *int:
		log.Debug(baseMsg + " int")

		tmp, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("could not parse value to int: %w", err)
		}

		*typedTo = tmp
	case *int8:
		log.Debug(baseMsg + " int8")

		tmp, err := strconv.ParseInt(value, 10, 8)
		if err != nil {
			return fmt.Errorf("could not parse value to int8: %w", err)
		}

		*typedTo = int8(tmp)
	case *int16:
		log.Debug(baseMsg + " int16")

		tmp, err := strconv.ParseInt(value, 10, 16)
		if err != nil {
			return fmt.Errorf("could not parse value to int16: %w", err)
		}

		*typedTo = int16(tmp)
	case *int32:
		log.Debug(baseMsg + " int32")

		tmp, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			return fmt.Errorf("could not parse value to int32: %w", err)
		}

		*typedTo = int32(tmp)
	case *int64:
		log.Debug(baseMsg + " int64")

		tmp, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("could not parse value to int64: %w", err)
		}

		*typedTo = tmp
	case *uint:
		log.Debug(baseMsg + " uint")

		tmp, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf("could not parse value to uint: %w", err)
		}

		*typedTo = uint(tmp)
	case *uint8:
		log.Debug(baseMsg + " uint8")

		tmp, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			return fmt.Errorf("could not parse value to uint8: %w", err)
		}

		*typedTo = uint8(tmp)
	case *uint16:
		log.Debug(baseMsg + " uint16")

		tmp, err := strconv.ParseUint(value, 10, 16)
		if err != nil {
			return fmt.Errorf("could not parse value to uint16: %w", err)
		}

		*typedTo = uint16(tmp)
	case *uint32:
		log.Debug(baseMsg + " uint32")

		tmp, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return fmt.Errorf("could not parse value to uint32: %w", err)
		}

		*typedTo = uint32(tmp)
	case *uint64:
		log.Debug(baseMsg + " uint64")

		tmp, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("could not parse value to uint64: %w", err)
		}

		*typedTo = tmp
	case *bool:
		log.Debug(baseMsg + " bool")

		tmp, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("could not parse value to bool: %w", err)
		}

		*typedTo = tmp
	case *float32:
		log.Debug(baseMsg + " float32")

		tmp, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return fmt.Errorf("could not parse value to float32: %w", err)
		}

		*typedTo = float32(tmp)
	case *float64:
		log.Debug(baseMsg + " float64")

		tmp, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("could not parse value to float64: %w", err)
		}

		*typedTo = tmp
	case *uuid.UUID:
		log.Debug(baseMsg + " uuid")

		tmp, err := uuid.Parse(value)
		if err != nil {
			return fmt.Errorf("could not parse value to uuid: %w", err)
		}

		*typedTo = tmp
	default:
		log.Debug("Unsupported type")

		return ErrUnsupportedType
	}

	log.Debug("Value correctly sat")

	return nil
}
