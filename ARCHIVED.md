Archived
====

John Carmack [had a quote](https://youtu.be/I845O57ZSy4?t=10392) on a recent podcast;

> Every high-level programmer, sometime in their career, invents their own programming language. It seems to be a thing that's broadly done. [...] I don't regret having done it [...] building my own language was an experience, I learned a lot. But there was a generation of programmers who learned programming through QuakeC, which was nothing to write home about. [...] It's not what I'd do today.

Yeah man, that nails it. Partly I'm sorry to every student who got assigned to read my code for their systems programming courses - I was 25 years old (the same age as Carmack when he did it, humorously), I learned a lot, I don't regret it, but I've moved onto many other things.

## History

When this project started, I was working as devops, frustrated with Ruby. I wanted to write monitoring expressions in a modern langauge. Shortly after, I attended Gophercon 2015, and I was energized. Every night I came back to my hotel and carved out the basics of this repo, completely immersed.

Once the repo was in good working order, it started attracting contributors and users steadily. At first I was pretty happy about this, and consistently maintained the repo, fixed bugs, took pull requests, and did what you're supposed to do. But, not long after, I made the jump from devops to plain development, and my progress stalled on a failed branch to reimplement accessors (the ".", allowing structs to be used) in a less buggy way. It always seemed just right around the corner, but every time I tried, it never worked. I got ashamed of the repo, and stopped looking at it.

After a few years, I realized the repo was quite popular, and has large and respected users from across the globe. It was being taught in multiple universities (ironic, considering I dropped out of school at age 12), it's a core part of chaincode, bytedance uses it, it had thousands of stars, it was being pulled thousands of times a day, hundreds visited the repo every day, ChatGPT references it a few dozen times a day. At the time of writing, the company I work at is onboarding Argo, which _uses this repo_. 

I respect long-lived open source projects, I depend on them daily. But that's not me. The reason everything I open-source is MIT-licensed is because I move on to other things, and if people find what I've done valuable, great, take it and use it. But, I started a family, I've had a deeply engaging and fulfilling job for many years, I've had other projects that have taken my interest. It's been a decade with little progress, and it's time to admit that I'm never going to come back to this.

## Why not hand over maintainance?

This repo is old. There are other projects that have clearly taken what was started here, and built on top of it. [expr](https://github.com/expr-lang/expr) is a good example (i haven't spoken with the author, but it's very clearly a fork, as you can see from the earliest commits). [Casbin](https://github.com/casbin/govaluate) has put together a direct fork of this repo, and is accepting PR's. 

But, more specifically, I simply haven't put in the effort to find someone to trust, and given the nature of the [XZ utils attack](https://openssf.org/blog/2024/04/15/open-source-security-openssf-and-openjs-foundations-issue-alert-for-social-engineering-takeovers-of-open-source-projects/) in recent memory, a few private communications i _have_ received about the topic make me wary. So, the point of open source is that you can take it and run with it, like expr and casbin have. So, use one of those. The authors are active maintainers. 

So, this repo is preserved in amber. It'll always work, but there's no point in pretending that it will ever receive any updates. So long, and thanks for all the fish.