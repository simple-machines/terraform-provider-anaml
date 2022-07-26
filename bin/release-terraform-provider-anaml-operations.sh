#!/usr/bin/env bash
set -euo pipefail

get_latest_tag() {
    local result
    result=$(git describe --tags --abbrev=0)
    echo "$result"
}



repo_has_operations_remote () {
    git remote get-url operations &> /dev/null && return
    false
}

is_terraform_provider_anaml_origin () {
    [[ "$(git remote get-url origin)" == "git@github.com:simple-machines/terraform-provider-anaml.git" ]] && return
    false
}

main () {
    if ! is_terraform_provider_anaml_origin; then
        echo "[ERROR] Need to run from terraform-provider-anaml repo"
        exit 1
    fi

    if [[ "$(git branch --show-current)" != "master" ]]; then
        echo "[ERROR] Please switch to master branch before continuing"
        exit 1
    fi


    if ! repo_has_operations_remote; then
        echo "[INFO] operations remote does not exist, creating, remote operations git@github.com:simple-machines/terraform-provider-anaml-operations.git"
        git add remote operations git@github.com:simple-machines/terraform-provider-anaml-operations.git
    fi

    echo "Preparing to sync terraform-provider-anaml to terraform-provider-anaml-operations for Terraform module publishing"
    while true; do
        read -r -p "Do you wish to continue? (Y/N): " answer
        case $answer in
            [Yy]* ) break;;
            [Nn]* ) exit 0;;
            * ) echo "Please answer Y or N.";;
        esac
    done

    echo "[INFO] fetching tags from origin"
    git fetch origin --tags

    local origin_commit_id
    local local_commit_id

    origin_commit_id="$(git log -n 1 --pretty=format:"%h" origin/master)"
    local_commit_id="$(git log -n 1 --pretty=format:"%h" origin/master)"

    if [[ ! "$origin_commit_id" == "$local_commit_id" ]]; then
        echo "[ERROR] Local master branch is not in sync with origin. Ensure you have pulled the latest master and pushed your local changes"
        exit 1
    fi

    local tag
    tag="$(get_latest_tag)"
    while true; do
        read -r -p "Do you want to push $tag to git@github.com:simple-machines/terraform-provider-anaml-operations.git? (Y/N): " answer
        case $answer in
            [Yy]* ) git push operations --tag "$tag"; break;;
            [Nn]* ) break;;
            * ) echo "Please answer Y or N.";;
        esac
    done

    while true; do
        read -r -p "Do you want to push master (commit $origin_commit_id) to git@github.com:simple-machines/terraform-provider-anaml-operations.git? (Y/N): " answer
        case $answer in
            [Yy]* ) git push operations master; break;;
            [Nn]* ) break;;
            * ) echo "Please answer Y or N.";;
        esac
    done

    echo "Done!"
}

main "$@"
