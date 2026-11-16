---
name: sync-spec
description: Bring spec.md back in line with the codebase after design changes. Checks recent code changes against the spec, updates only sections that are now wrong or missing a new design decision, and adds decision-log rows. Use when the user runs /sync-spec or asks to update the spec to match the code.
disable-model-invocation: true
---

# Sync spec.md with the code

`spec.md` documents the **design** of senior_intern_bot: behavior, invariants,
interfaces, data model, sources, and the reasons behind each decision. The
code is the source of truth for what the system does. The user is the only
source of truth for *why*.

## 1. Work out what changed

1. Find the last commit that touched the spec:
   `git log -1 --format='%h %ad %s' --date=short -- spec.md`
2. List the code changes since then, including uncommitted work:
   `git log --oneline <that-commit>..HEAD` and `git diff <that-commit> --stat`
3. If the user named a branch, commit range, or feature, use that instead.
4. Check which branch holds the code being documented. If it isn't the
   current branch, say so before editing, because the spec would then
   describe code that isn't on this branch.

## 2. Sort changes into design changes and refactors

Read the actual diffs, not just the commit messages.

**Needs a spec update** if any of these is true:
- A statement in the spec is now false (a number, a behavior, an ordering,
  a guarantee, a field, an endpoint).
- The change makes a new design decision: a new invariant, a new failure
  mode or how it's handled, a changed trade-off, a new source, a new stage,
  a changed schema, a changed timeout or interval.
- A decision the spec lists as deferred or planned has been built.

**Leave the spec alone** for refactors: renames, moved code, extracted
helpers, readability changes, new tests, lint fixes. If behavior and design
are unchanged, the spec is unchanged. When unsure, list the change under
"skipped" in the report rather than editing.

## 3. Verify every claim against the code

For each section you will touch, read the current code it describes and
check every concrete claim: numbers, limits, call order, what is persisted,
what is retried. Do not rely on commit messages, earlier conversation, or
memory. Counts (e.g. number of watched companies) come from the code.

## 4. Edit

- Change only what is wrong or missing. Don't reword, reorder, or restyle
  sections that are still accurate.
- Match the existing style: numbered sections, short paragraphs that explain
  *why*, tables for structured facts, backticks for identifiers.
- New topics go in a new numbered subsection next to the related ones. Keep
  cross-references (`§5`, `§8.2`) correct.
- Historical measurements that explain a past decision stay, but say they
  were true at the time if the code has since changed.
- Add a row to the decision log (§9) for each new design decision:
  `| YYYY-MM | <decision> | <reason> |`, using the commit date's month.

**Never invent a reason or a number.** If the code shows *what* changed but
nothing says *why* (e.g. a timeout raised with no explanation), or a figure
would need measuring, leave it out of the spec and ask the user in the
report. Do not fill gaps with a plausible-sounding rationale.

If `roadmap.md` exists and tracks the same item, update its status too.

Do not commit. The user reviews and commits.

## 5. Report back

Keep it short:
- **Updated:** each section changed, and what it now says, in one line.
- **Skipped:** changes judged to be refactors, so the user can disagree.
- **Questions:** missing reasons or numbers needed to finish the update.
- **Branch note:** if the documented code isn't on the current branch.
