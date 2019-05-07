#!/bin/bash

if [[ $# -lt 2 ]]
then
  echo "Use: "$0" <url> <directory_path>"
  echo "Clone webpage & convert to UTF-8"
  exit
fi

cd site/

# Download webpage
wget -nH --user-agent="Mozilla/5.0 (Windows NT 5.2; rv:2.0.1) Gecko/20100101 Firefox/4.0.1" -p -np --restrict-file-names=nocontrol -P $2 $1

# Remove query string from extension
for i in `find $2 -type f -name "*\?*"`; do
    mv $i `echo $i | cut -d? -f1`; 
done

cd $2

# Convert ot UTF-8
for i in `find . -type f -exec grep -Iq . {} \; -print`; do
  if [[ ! -f $i ]]; then # Only convert text files
    continue
  fi

  echo "Converting $i to UTF-8"
  # Generate temp file to avoid Bus error
  iconv -f `uchardet $i` -t utf-8 $i -o $i.tmp
  mv $i.tmp $i

  # Inject script that removes hidden fields from forms
  sed -i "s/\(<\/body>\)/<script src=\"\/inject.js\"><\/script>\1/gi" $i
done

cd ../

# Copy inject script to directory
cp ../inject.js "$2/"
