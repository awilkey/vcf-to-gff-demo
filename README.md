# vcf-to-gff
Conversion between vcf and gff using [bio-format-tools-go](https://github.com/awilkey/bio-format-tools-go)
defaults to reading from stdin and outputting to stdout

flags (optional)

`-s <file>`  input file (default:stdin)

`-d <file>`  output file (default:stdout)

`-c` count numner of lines in input vcf (used for padding when naming)

`-l #` number of places for left-based zero padding (default: 7)

This tool converts from a snpeff vcf format:

`Chr01   1012    .       A       T       2158.27 PASS    ANN=T|intergenic_region|MODIFIER|CHR_START-Glyma.01G000100|CHR_START-Glyma.01G000100.Wm82.a2.v1|intergenic_region|CHR_START-Glyma.01G000100.Wm82.a2.v1|||n.1012A>T||||||        GT:AD:DP:GQ:PL  [...]`          

to a gff of format:

`Gm01    SNP     HapMap  1012    1013    .       +       .       Name=A03.000008;Effect=A,T,SNP,intergenic_region,MODIFIER,Glyma.01G000100`
