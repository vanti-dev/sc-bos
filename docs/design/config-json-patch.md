# User editable config: JSON Patch

Here we discuss what a user-editable config solution would look like if we used JSON Patch to record user edits.

Here's the general gist:
We split config into two parts, system config and user config.
System config is supplied by us (or other contractor) and typically is generated via tools and uploaded in bulk.
When you click "reset" you are resetting back to the system config.
User config is config applied via forms in the Ops UI, for example changing the timeout on a PIR sensor.
This config is recorded as a JSON Patch (or collection of JSON Patch) documents and applied to the system config to get
the final configuration values.

JSON Patch is a standard described by [RFC 6902](https://tools.ietf.org/html/rfc6902) and encodes a sequence of
operations to modify a JSON structure.
We'd probably want to split these patches into "sessions" adding our own metadata to the patch, like who did it and
when, maybe a comment or something like that.

Here's an example JSON Patch document:

```json
[
  {"op": "test", "path": "/a/b/c", "value": "foo"},
  {"op": "remove", "path": "/a/b/c"},
  {"op": "add", "path": "/a/b/c", "value": ["foo", "bar"]},
  {"op": "replace", "path": "/a/b/c", "value": 42},
  {"op": "move", "from": "/a/b/c", "path": "/a/b/d"},
  {"op": "copy", "from": "/a/b/d", "path": "/a/b/e"}
]
```

We'd likely end up with user edit files more like this

```json
{
  "owner": {"title": "Matt"},
  "description": "Replace faulty PIRs and light fixture IDs in DALI driver",
  "changeTime": "2021-01-01T12:00:00Z",
  "patch": [
    {"op": "replace", "path": "/drivers/dali/0/lights/0/address", "value": 42},
    {"op": "replace", "path": "/drivers/dali/0/lights/0/metadata/id", "value": "new-id"},
    {"op": "replace", "path": "/drivers/dali/0/lights/1/address", "value": 43},
    {"op": "replace", "path": "/drivers/dali/0/lights/1/metadata/id", "value": "new-id"}
  ]
}
```

I'd also expect that we'd maintain a document containing the generated config (system + patches) to improve boot
performance.

## Concerns

### Arrays that are really sets

We have quite a few config structures that are represented as arrays, but are really sets.
That is to say the order of these arrays doesn't matter.

In these cases we try to maintain idempotency by sorting the arrays using some key before writing them to disk (and
committing to Git).
This works well when we are looking at git diffs but isn't great when using JSON Patch, which
uses [JSON Path](https://tools.ietf.org/html/rfc6902#section-4) as adding a new item to the start of an array at the
system level would break all paths in user code, which now are referring to the wrong array item.

We could modify existing user edits to make their indexes match the new array items.
This assumes we can tell that a new set of system config represents an array insert.
To make that work we'd need to define keys for any array that is really a set and encode that into our code somewhere.

### Performance

As mentioned before, it'll be likely that we maintain a cache config that contains all user edits applied to the system
config.
During runtime this will be incrementally updated without any performance impact.

Updating the system config, would invalidate this cache and require a full re-apply of all user edits.
This could take arbitrarily long, and introduces a process that by design would take more time the longer the system is
in use.

Another performance impact might come from config that is regularly "toggled".
With a simple design we'd create new patch documents each time the property is changed, then another when changed back.
For properties that end up being changed regularly we'd end up with a lot of files that only grow over time.

We could potentially solve both issues by culling user edits, removing overlapping changes.
If two user edits both change the same property, only the latest edit would be kept.

### Future proofing

Over time I hope and expect the distribution of configuration between system and user to shift.
Right now 100% of config is system config - there's no way to enter user config, once this feature is implemented that
ratio will change, but not by much as there won't be many things you can configure via the Ops UI.
However, as we add more and more user configurable forms, and as we focus more on commissioners being able to setup an
SC BOS system, I hope that the config will eventually become 100% user config and the system config just exists as the
default out-of-box experience.

This design is based on the current situation and isn't designing for the future situation.
In other words, we're optimising for the case we think will go away over time.

This is not to say we couldn't use some form of this proposal for ad-hoc user changes.
Maybe there's a grace period where we collect user edits as patches then there's a "commit" action you perform to "save"
the collected config as something we want to commit to.
