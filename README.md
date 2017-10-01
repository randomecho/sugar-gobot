# Sugar Go bot

Tool to quickly create Account and Contact records in a Sugar instance
via the REST API. Each Account is linked with a Contact.

## Usage

To create three records (six in total of Accounts, Contacts) and leave in the system:

    sugar-gobot -num=3

To create 10 records (20 in total of Accounts, Contacts) and delete them right after:

    sugar-gobot -num=10 -delete=true


## Setup

Create a copy of **config/app.yml** from the example **config/app.example.yml**

    $ go get
    $ go install
    $ sugar-gobot


### License

[Released under the MIT license](/LICENSE).
