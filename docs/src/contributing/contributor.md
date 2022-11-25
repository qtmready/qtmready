# Considerations for Contributor

The responsibility to find the best solution and implement it lies with the contributor.

Before asking for a review, the contributor should be confident that their code

- actually solves the problem it was meant to solve.
- does so in the most appropriate way.
- satisfies all requirements.
- are no remaining bugs, logical problems, uncovered edge cases, or known vulnerabilities.

The best way to do this, and to avoid unnecessary back-and-forth with reviewers, is to perform a self-review of your own merge request.

## Each PR must have a good description

From reading the description, the reviewer should be able to understand what the code is meant to do. This has to be true even if there is a shortcut story or a requirements page.

PRs without a description will never go through because the reviewer of a PR without a description has 2 jobs:

- Understanding what the code does from reading the code changes in the PR
- Trying to decide if the code does what it’s supposed to do based on the reviewer’s general views on the universe and the current geopolitical climate?

If there is no statement that clearly explains what the code is meant to do, reviewing it for correctness is not possible. You simply don’t know what is correct. You operate on assumptions of what you think is correct.

## PR must have sufficient unit test and integration test coverage

Ideally, the reviewer has an easy way to see the list of test cases covered.

If test coverage already exists, or for whatever reason test coverage is not complete (work split off into another ticket, dependencies), _PR_ description should mention this.

**_Giving good PR reviews are hard and time-consuming_**. That’s because once you understand what the code is supposed to do, you have to stop reviewing and jot down what test coverage you’d want for this. You make your own test cases list and decide what kind they are: unit test/integration/end-to-end.
You also have to think if there are possible production implications, and if any testing on live data is needed pre-deploy.

## For a bug fix, test-cases must be updated

For bug fixes, a test case must be written so that, should the bug fix be reverted, this test would fail.

Say there’s a typo bug where you went != instead of ==. It was 10pm and you shouldn’t have been coding anyway, but you felt like going the extra mile that night. So you open a bug fix PR and change != to ==. Some kind soul, also working, approves it, and off it goes to the master branch.

Three days later your buddy Duke is finally ready with his PR that got out of control: he thought it was a small feature and didn’t break it into smaller tickets, but that didn’t turn out to be the case and now he has a 2000 line PR he is begging his team to review.

Duke merges from upstream, the code conflicts where you swapped == in for !=. Duke is tired, he doesn’t notice what’s up, != is back.

Since there’s no test to catch it, the bug’s back.

This is a contrived example, but if you pay attention over longer periods — this happens all the time. Lay out the bugs over 3 years, you’ll see those not trapped by a test reappearing.
