#!/bin/sh

# Exit on error
set -e
# Remove all the test files when we're done
trap owltest_teardown EXIT

# Check PWD
if [ "$(basename $PWD)" != owl ]; then
    echo 'Run this from project top level, i.e. "owl"'
    exit 1
fi

export OWL_START=$PWD
export OWL_TESTDATA=$PWD/testdata

# Runes that I don't want to write
export mouse_emoji=$'\U0001F401'
export germany_emoji=$'\U0001F1E9\U0001F1EA'
export bird_emoji=$'\U0001f426'
# 猫头鹰
export chinese_simplfied_owl=$'\u732B\u5934\u9E70'

owltest_setup() {
    mkdir "$OWL_TESTDATA"
    cd "$OWL_TESTDATA"
    mkdir -p 'birds:/pengu?ns' 'ascii_dir'
    touch \
        'no-change-needed' \
        'So * many* *** *c*veats' \
        '::>?' \
        '<>??' \
        'birds:/tawny: the game' \
        'birds:/pengu?ns/pingu?' \
        'birds:/pengu?ns/rockhopper' \
        'birds:/pengu?ns/rockhopper?' \
        "Mouse? ${mouse_emoji}" \
        "Deustchland, Deutschland! ${germany_emoji}" \
        "${chinese_simplfied_owl}" \
        $'Invalid\x80UTF-8' \
        $'fich\u00E9-en-fran\u00E7ais' \
        $'No\u0338rse\u0301' \
        'Question? Why' \
        'QUESTION: WHY'
}

owltest_teardown() {
    cd "$OWL_START"
    rm -r "$OWL_TESTDATA"
}

owltest_reset() {
    owltest_teardown
    owltest_setup
}

owl() {
    $OWL_START/owl "$@"
}

unicode_range() {
    local low=$(printf %d $1)
    local high=$(printf %d $2)
    for i in $(seq $1 $2); do
        printf "0d%d_0x%x -> \U$(printf '%08x' $i)\n" $i $i
    done |
        column
}

# Help message
cat <<-EOF
Testing script; use this to:
  - use the owl executable in this directory as just 'owl'
  - have a set of files and directories to test out owl with
This will be a lot easier if you have 'tree' installed; see README.md
for details.

You will now be in ${OWL_TESTDATA}; cd out and delete this directory to 
clean up after a test. Or run owltest_reset to do the same in-place.
EOF

export -f owltest_setup owltest_reset owltest_teardown owl unicode_range
owltest_setup
bash -i
echo '+++ Test data directory cleared +++'
