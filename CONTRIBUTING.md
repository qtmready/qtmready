# Contributing to this repo

Oh there you are, we have been expecting you. Let us get right to it then.

## Teamwork

The most important thing is to value your team-mates and their time. Since git is where the coding happens, we have formulated a number of guidelines to effectively contribute as a team. The general guidelines are

- Be kind.
- Accept that many programming decisions are opinions. Discuss tradeoffs, which you prefer, and reach a resolution quickly.
- Ask questions; don’t make demands. (“What do you think about naming this :user_id?”)
- Ask for clarification. (“I didn’t understand. Can you clarify?”)
- Avoid selective ownership of code. (“mine”, “not mine”, “yours”)
- Avoid using terms that could be seen as referring to personal traits. (“dumb”, “stupid”). Assume everyone is attractive, intelligent, and well-meaning.
- Be explicit. Remember people don’t always understand your intentions online.
- Be humble. (“I’m not sure - let’s look it up.”)
- Don’t use hyperbole. (“always”, “never”, “endlessly”, “nothing”)
- Be careful about the use of sarcasm. Everything we do is public; what seems like good-natured ribbing to you and a long-time colleague might come off as mean and unwelcoming to a person new to the project.
- Consider one-on-one chats or video calls if there are too many “I didn’t understand” or “Alternative solution:” comments. Post a follow-up comment summarizing one-on-one discussion.
- If you ask a question to a specific person, always start the comment by mentioning them; this will ensure they see it if their notification level is set to “mentioned” and other people will understand they don’t have to respond.

## Commit Message

Our entire workflow revolves around [shortcut](https://app.shortcut.com/ctrplane) and our git messages are designed after shortcut stories.
Suppose we have a user story with title `Fix MQTT message delivery` and shortcut gives it a number 99. When we are done with the fix, the git message will be designed like

```docs
[sc99]: fix for MQTT message delivery

BREAKING CHANGES (Optional):
if there are any breaking changes.

TODO (Optional):
if there is anything left to do and if we are waiting for some other roadblock.
```

## Merge/Pull Requests

Avoid large pull requests. Find ways to break them up into smaller ones. Small PRs can be:

- Easier to review. The intent is clearer and mistakes are more obvious.
- Less likely to conflict with others work
- Compatible w/large new features - e.g behind a flag

TIP: The easiest way to do is, have a shortcut story refer to a very small use case. If you feel it is growing, convert the collection of use cases into an epic.

### Code Reviews

#### Considerations for the contributor

The responsibility to find the best solution and implement it lies with the contributor.

Before asking for a review, the contributor should be confident that their code

- actually solves the problem it was meant to solve.
- does so in the most appropriate way.
- satisfies all requirements.
- are no remaining bugs, logical problems, uncovered edge cases, or known vulnerabilities.

The best way to do this, and to avoid unnecessary back-and-forth with reviewers, is to perform a self-review of your own merge request

##### Each PR must have a good description

From reading the description, the reviewer should be able to understand what the code is meant to do. This has to be true even if there is a shortcut story or a requirements page.

PRs without a description will never go through because the reviewer of a PR without a description has 2 jobs:

- Understanding what the code does from reading the code changes in the PR
- Trying to decide if the code does what it’s supposed to do based on… the reviewer’s general views on the universe and the current geopolitical climate?

If there is no statement that clearly explains what the code is meant to do, reviewing it for correctness is not possible. You simply don’t know what is correct. You operate on assumptions of what you think is correct.

##### PR must have sufficient unit test and integration test coverage

Ideally, the reviewer has an easy way to see the list of test cases covered.

If test coverage already exists, or for whatever reason test coverage is not complete (work split off into another ticket, dependencies), _PR_ description should mention this.

**_Giving good PR reviews are hard and time-consuming_**. That’s because once you understand what the code is supposed to do, you have to stop reviewing and jot down what test coverage you’d want for this. You make your own test cases list and decide what kind they are: unit test/integration/end-to-end.
You also have to think if there are possible production implications, and if any testing on live data is needed pre-deploy.

##### For a bug fix, test-cases must be updated

For bug fixes, a test case must be written so that, should the bug fix be reverted, this test would fail.

Say there’s a typo bug where you went != instead of ==. It was 10pm and you shouldn’t have been coding anyway, but you felt like going the extra mile that night. So you open a bug fix PR and change != to ==. Some kind soul, also working, approves it, and off it goes to the master branch.

Three days later your buddy Duke is finally ready with his PR that got out of control: he thought it was a small feature and didn’t break it into smaller tickets, but that didn’t turn out to be the case and now he has a 2000 line PR he is begging his team to review.

Duke merges from upstream, the code conflicts where you swapped == in for !=. Duke is tired, he doesn’t notice what’s up, != is back.

Since there’s no test to catch it, the bug’s back.

This is a contrived example, but if you pay attention over longer periods — this happens all the time. Lay out the bugs over 3 years, you’ll see those not trapped by a test reappearing.

#### Considerations for the reviewer

The code is

- Well designed
- Readable by others
- Doing what the author intended
- No more complex than needed
- Not degrading system code health
- Commented with the why vs. what
- Appropriately tested
- Sufficiently documented

#### How to provide feedback on code

- Assume competence (author may have a different context)
- Provide specific, actionable feedback
- Focus on the code (avoid personal criticism)
- Mark nitpicks & optional comments
- Provide rationale (e.g refs to a style guide, link to a blog, etc)

### Merging responsibility

At least two reviews are required before one can merge back. Manager approval doesn't count. Yes, that is right! Maybe you a manager and this made you think, “**_WHAT? I am a manager, and an amazing coder, and my review totally counts, how dare you!!_**”. But you have to consider how much time your manager spends coding. Is it more than 50%? Then they are also a developer, so okay, let’s count them. At the end of the day, it is the developers who have to live with the code that’s being written. Therefore developers get the final say. What we aim to prevent is a manager-developer infinite approval cycle where a manager keeps approving one developer’s PRs even though there are 3 other developers on the team who haven’t had a chance to look at it yet.

## Project Planning & Management

The project in its entirety runs around [shortcut](https://shortcut.io).

### Stories

A small actionable use case that can be fit into a single Pull/Merge Request. See Merge Requests howto above. To create a story

- Provide complete information on what the user intends to achieve.
- Provide complete success criteria.
- Update the story if we discover something, a bug, a related action item, so that others can learn from it.
- while working, associate your commit with the story it is impacting. see commit messages above.

### Epics

Epic is a collection of user stories that are related by functionality. For example, an epic can be _tag sharing_ with stories as _facebook sharing_, _linkedin share_, _twitter share_, _facebook developer applicaiton_, _twitter developer applicaiton_, etc

### Milestones

Milestones are larger pieces of work with a collection of epics and individual stories.
