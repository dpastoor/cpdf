cpdf
====

cpdf is a super lightweight utility to quickly combine pdf files in a directory together.

```
cpdf <output_file_name.pdf> 
```

will combine all files into the output file, minus the output file name, if it exists

For targetted combination, regex matching can be used

```
cpdf <output_file_name.pdf> <regex_match>
```

check the version with

```
cpdf -v
```