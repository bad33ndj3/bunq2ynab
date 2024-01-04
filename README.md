# YNAB to BUNQ

![image](https://github.com/bad33ndj3/bunq2ynab/assets/9072952/2a61635f-db5e-43a2-825a-8a230cd3bebc)


This script will sync BUNQ transactions to YNAB.

## WIP

This script is still in development since I just started testing it.

#### TODO:

- [x] Add support for multiple budgets
- [x] Optimize API calls
- [ ] Fix internal transfers
- [ ] Add Joint account support
- [ ] Add more tests
- [ ] Add more documentation
- [ ] Add CI/CD

## Installation

1. Clone this repository
2. Setup config file `cp example.config.yml config.yml`
3. Fill in the config file with your own data
    - bunq_token can be found in the bunq app
    - ynab_token can be found in the YNAB settings
    - accounts is a list of accounts to sync
        - bunq_account_name is the name of the account in bunq
        - ynab_budget_name is the name of the budget in YNAB (Top level)
        - ynab_account_name is the name of the bank account in YNAB
4. Run `make run`
5. Wait for the script to finish

## Similar projects
- [ynab](https://support.ynab.com/en_us/direct-import-in-the-uk-and-eu-an-overview-Syae1z_A9) Last year YNAB added support for direct import in the UK and EU.  This is a great alternative if your bank is supported.
- [bunq2ynab](https://github.com/wesselt/bunq2ynab) Python script to import transactions from bunq bank to YNAB.  Supports listening to messages from bunq so your payments show up in YNAB seconds after you pay.
- [syncforynab](https://syncforynab.com/) Although this is a paid service, it's a great alternative if you don't want to run your own script.
- [awesome-ynab](https://github.com/scottrobertson/awesome-ynab) A curated list of awesome things related to You Need A Budget (YNAB).
