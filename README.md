# Tag Issues Which have not moved in project board

The goal of this action is to label issues which have been in a column in a project for a certain amount of time.

This action can be helpful if you use project boards to manage stages of a project or a team and want to automatically label issues if they have been sitting in a column too long.

## Inputs

* PROJECT_ID - id of the project you want to set up this action for
* TIME_IN_MINUTES - issues which have stayed in a project board column for this time period atleast get flagged
* LABEL_NAME - name of label to be added to issues exceeding time in column

## Sample workflow 

```
name: Example
on: 
  workflow_dispatch:
  project_card:
    types: [created, moved]
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/cache@v2
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    - uses: actions/setup-go@v2
      with:
        go-version: '1.15.8'
    - name: Example
      uses: isethi/issue-timer@main
      with:
        PROJECT_ID: 1
        TIME_IN_MINUTES: 10
        LABEL_NAME: review_time
```
