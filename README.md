# Cronex

Cronex is a very basic cron expression parser. It supports five fields and a command separated by
whitespace. 

Table 1 shows the fields in the expected order.


## Table 1 - Cron Expresssions Allowed Fields and Values
| Name | Allowed Values | Features |
|------|:----------:|:-----------------------:|
|Minutes| 0-59 | , - */|
|Hours| 0-23 | , - */|
|Day of month| 1-31 | , - */|
|Month| 1-12 | , - */|
|Day of week| 1-7 | , - */|

## Getting started

```
git clone https://github.com/n0block/cronex

cd cronex

go build

cronex "*/15 5-21 1,15 * 1-3 /usr/bin/find"
```