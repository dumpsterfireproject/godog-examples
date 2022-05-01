Feature: Account Maintenance

As an account holder, I need to be able to depoit and withdraw money from my account,
and I expect the bank to be able to track the balance of my account accurately at all times.

Background: Setup new account
Given I have a new account

# Rule: ???

@foo @bar@baz @category(integration)
@dogg
Scenario: New account
Given I have a new account
 Then the account balance must be 0.00 USD

@category(integration)
Scenario: Deposit money into account
Given I have an account with 0.00 USD
 When I deposit 5.00 USD
 Then the account balance must be 5.00 USD

@bar
Scenario: Withdraw money from account
Given I have an account with 11.00 USD
 When I withdraw 5.00 USD
 Then the account balance must be 6.00 USD

@priority:high
Scenario: Attempt to overdraw account
Given I have an account with 11.00 USD
 When I try to withdraw 50.00 USD
 Then the transaction should error

Scenario: Concurrent Deposits and withdrawals
Given I have an account with 100.00 USD
 When I process the following transations:
|type       |dollars|
|deposit    |5      |
|deposit    |10     |
|withdrawal |2      |
|withdrawal |20     |
 Then the account balance must be 93.00 USD

@concurrency
Scenario Outline: Opening accounts in different currencies
Given I have an account with 100.00 <currency>
 Then the account balance must convert to <dollars> USD

Examples:
|currency|dollars|
|CAD     |80     |
|CNY     |16     |
|EUR     |108    |

@wip
Examples:
|currency|dollars|
|CAD     |80     |
|CNY     |16     |
|EUR     |108    |

@issue#952
Scenario: Remittance address
Given I have a new account
 Then the remittance address must be
"""
742 Evergreen Terrace
Springfield, OR
"""
