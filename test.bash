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
    $OWL_START/owl $@
}

unicode_rune() {
    printf "\U$(printf '%08s' $1)\n"
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

export -f owltest_setup owltest_reset owltest_teardown owl unicode_rune
owltest_setup
bash -i
