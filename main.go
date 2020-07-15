// package gcvit provides handlers for the gcvit server
// API routes
package main

import (
	"bufio"
	"bytes"
	//"encoding/json"
	"compress/gzip"
	"flag"
	"fmt"
	"github.com/awilkey/bio-format-tools-go/gff"
	"github.com/awilkey/bio-format-tools-go/vcf"
	"io"
	"log"
	"strconv"
	"strings"
	//"time"
	"os"
)

func main() {
	// setup flags
	var input = flag.String("s","","Input file if not using stdin")
	var output = flag.String("d","","Output file if not using stdout")
	var padding = flag.Int("l",7,"estimated number of lines")
	var count = flag.Bool("c",false,"count lines before processing")

	// parse passed flags
	flag.Parse()
	pad := *padding
	var err error

	// Count lines if requested
	if *count {
		pad,err = calcPadding(*input)
		if err != nil {
			log.Fatal("error reading vcf file: ", err)
		}
	}

	//Open vcf file for parsing
	vcfReader, err := newVcfReader(*input)
	if err != nil {
		log.Fatal("error reading vcf file: ", err)
	}

	//Print gff
	var b bytes.Buffer
	writer, err := gff.NewWriter(&b)
	if err != nil {
                log.Fatal("Error: Problem opening gff writer: %s", err)
        }

	out := os.Stdout
	if *output != "" {
		out, err =  os.Create(*output)
		if err != nil {
			log.Fatal("error opening gff file: ", err)
		}
		defer out.Close()
	}

	var feat *vcf.Feature // line of vcf file
	var readErr error
	var line int
	line = 1
	// Read VCF file and format/print gff
	for readErr == nil {
		feat, readErr = vcfReader.Read()
		if feat != nil {
			arr := strings.Split(feat.Info["ANN"], "|")
			arr[3] = strings.Replace(arr[3], "CHR_START-","",-1)
			gffLine := gff.Feature{
				Seqid:      strings.Replace(feat.Chrom, "Chr", "Gm", -1),
				Source:     "HapMap",
				Type:       "SNP",
				Start:      feat.Pos,
				End:        feat.Pos +1,
				Score:      gff.MissingScoreField,
				Strand:     "+",
				Phase:      gff.MissingPhaseField,
				Attributes: map[string]string{"Name": fmt.Sprintf("%s ;", LPad(strconv.Itoa(line), pad)), "Effect": fmt.Sprintf("%s,%s,SNP,%s,%s,%s", feat.Ref, strings.Join(feat.Alt,","),arr[1],arr[2],arr[3])},

			}
			line++
			//write line to buffer
			writer.WriteFeature(&gffLine)
			//write line to output
			fmt.Fprintf(out, "%s", b.String())
			//clear buffer to save on RAM
			b.Reset()
		}
        }
}



// Open VCF file as a vcfReader
func newVcfReader(loc string) (*vcf.Reader, error) {
	if loc == "" {
		return vcf.NewReader(os.Stdin)
	} else {
		gz := strings.Contains(loc, ".gz")
		f,err := os.Open(loc)
		if err != nil {
			return nil, fmt.Errorf("problem with os.Open: %s",err)
		}
		if gz {
			contents, err := gzip.NewReader(f)
			if err != nil {
				return nil, fmt.Errorf("problem with gzip.NewReader(): %s",err)
			}
			return vcf.NewReader(contents)
		} else {
			return vcf.NewReader(f)
		}
	}
}

//Left-pad with 0
func LPad(s string, pLen int) string {
	var padStr = strings.Repeat("0",pLen) + s;
	return padStr[(len(padStr) - pLen):]
}

//Count lines in file for padding
func calcPadding(loc string) (int, error){
	var count int
	const lineBreak = '\n'
	var r io.Reader
	buf := make([]byte, bufio.MaxScanTokenSize)
	if loc == "" {
		r = os.Stdin
	} else {
		gz := strings.Contains(loc, ".gz")
		f,err := os.Open(loc)
		if err != nil {
			return 0, fmt.Errorf("problem with os.Open: %s",err)
		}
		if gz {
			r, err = gzip.NewReader(f)
			if err != nil {
				return 0, fmt.Errorf("problem with gzip.NewReader(): %s",err)
			}
		} else {
			r = f
		}
	}

	for {
		bufferSize, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}

		var buffPosition int
		for {
			i := bytes.IndexByte(buf[buffPosition:], lineBreak)
			j := bytes.IndexByte(buf[buffPosition:], '#')
			//break if no more linebreaks
			if i == -1 || bufferSize == buffPosition {
				break
			}
			buffPosition += i + 1
			//skip header lines
			if j == 0 {
				continue
			}
			count++
		}
		//break if eof
		if err == io.EOF {
			break
		}
	}

	fmt.Fprintf(os.Stderr, "number of lines: %d \n", count)

	pad := 0
	for count != 0 {
		count /= 10
		pad++
	}

	return pad, nil
}
