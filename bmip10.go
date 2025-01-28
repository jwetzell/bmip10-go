package bmip10

type Config struct {
	SampleStates       uint32
	CodedStates        uint32
	LossyCodeWidth     uint32
	LossyRounding      uint32
	LosslessCodes      uint32
	LossyCodes         uint32
	CodeTables         uint32
	DefaultTable       uint32
	LowestTableThresh  uint32
	HighestTableThresh uint32
	LossyFlag          uint32
}

func SetupBMIP10(sampleBits uint32, codedBits uint32) *Config {
	config := Config{}
	config.SampleStates = 1 << sampleBits
	config.CodedStates = 1 << codedBits
	config.LossyCodeWidth = (1 << (sampleBits - codedBits + 1)) - 1
	config.LossyRounding = config.LossyCodeWidth / 2
	config.LosslessCodes = config.CodedStates / 2
	config.LossyCodes = config.CodedStates / 2
	config.CodeTables = 1 + config.CodedStates/2
	config.DefaultTable = config.CodeTables / 2
	config.LowestTableThresh = config.LosslessCodes / 2
	config.HighestTableThresh = config.SampleStates - (config.LosslessCodes / 2)
	config.LossyFlag = (1 << codedBits) / 2

	return &config
}

// Use the decoded sample to select the next table to use
func NextTable(config *Config, sample int32) int32 {
	if sample < int32(config.LowestTableThresh) {
		return 0
	} else if sample >= int32(config.HighestTableThresh) {
		return int32(config.LossyCodes)
	} else {
		return (sample - int32(config.LowestTableThresh) + int32(config.LossyRounding)) / int32(config.LossyCodeWidth)
	}
}

// For a given sample and code table calculate the code word
func EncodeSample(config *Config, table int32, sample int32) int32 {
	lossless_low := table * int32(config.LossyCodeWidth)
	lossless_high := lossless_low + int32(config.LosslessCodes)

	if sample >= lossless_low && sample < lossless_high {
		return sample - lossless_low
	} else if sample < lossless_low {

		return int32(config.LossyFlag) | (sample / int32(config.LossyCodeWidth))
	} else {
		return int32(config.LossyFlag) | ((sample-lossless_high)/int32(config.LossyCodeWidth) + table)
	}
}

// For a given code word and code table decode the sample
func DecodeSample(config *Config, table int32, code_word int32) int32 {
	index := code_word & ^int32(config.LossyFlag)
	if (code_word & int32(config.LossyFlag)) != 0 {
		if index < table {
			return index*int32(config.LossyCodeWidth) + int32(config.LossyRounding)
		} else {
			return index*int32(config.LossyCodeWidth) + int32(config.LossyRounding) + int32(config.LosslessCodes)
		}
	} else {

		return table*int32(config.LossyCodeWidth) + index
	}
}

type Encoder struct {
	Config *Config
	Table  int32
}

func NewEncoder(sampleBits uint32, codedBits uint32) *Encoder {
	config := SetupBMIP10(sampleBits, codedBits)
	return &Encoder{
		Config: config,
		Table:  int32(config.DefaultTable),
	}
}

func (e *Encoder) Encode(sample int32) int32 {
	code_word := EncodeSample(e.Config, e.Table, sample)
	decoded := DecodeSample(e.Config, e.Table, code_word)
	e.Table = NextTable(e.Config, decoded)
	return code_word
}

func (e *Encoder) Reset() {
	e.Table = int32(e.Config.DefaultTable)
}

type Decoder struct {
	Config *Config
	Table  int32
}

func NewDecoder(sampleBits uint32, codedBits uint32) *Decoder {
	config := SetupBMIP10(sampleBits, codedBits)
	return &Decoder{
		Config: config,
		Table:  int32(config.DefaultTable),
	}
}

func (d *Decoder) Decode(code_word int32) int32 {
	decoded := DecodeSample(d.Config, d.Table, code_word)
	d.Table = NextTable(d.Config, decoded)
	return decoded
}

func (d *Decoder) Reset() {
	d.Table = int32(d.Config.DefaultTable)
}
