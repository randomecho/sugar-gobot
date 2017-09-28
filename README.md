# Sugar Go bot

Tool to quickly create Account and Contact records in a Sugar instance
via the REST API. Each Account is linked with a Contact.

## Usage

To create 10 records (20 in total of Accounts, Contacts) and delete them right after:

    sugar-gobot -num=10 -delete=true

To create three records (six in total of Accounts, Contacts) and leave in the system:

    sugar-gobot -num=3


### License

[Released under the MIT license](/LICENSE).
