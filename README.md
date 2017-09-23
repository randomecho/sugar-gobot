# Sugar Go bot

Connection class to Sugar CRM API.

## Usage

To create 10 records (20 in total of Accounts, Contacts) and delete them right after:

    sugar-gobot -num=10

To create three records (six in total of Accounts, Contacts) and leave in the system:

    sugar-gobot -num=3 -delete=false
