# Thoughts about config v2 for SC BOS

The current config system (appconf and sysconf) is a bit fragile. While we do have a way to store user changes, the
process isn't very robust and lack a lot of features we'd want from a production config system. Some alternatives have
been designed - notably the Publication trait - but these were never integrated into sc-bos and don't satisfy all the
possible requirements we'd want.

## Wish list

1. Support both bulk/tool generated config and manual user edits
    1. Be able to see clearly what a user changed verses the tool generated config
    2. Be able to re-run the tool and not lose the user changes
    3. Be able to run the tool while the user is making changes - maybe something like the Refresh button you get in
       GitHub PRs when changes have been pushed while you're reviewing.
2. Support for editing config on other nodes - most notably editing config via edge gateways for an area controller
3. Maybe support both push and pull models for distribution.
    - push has be benefit that you know something has been done but is also difficult to scale across infrastructures
    - pull is less deterministic but also more likely to work in distributed environments like over the web - the web is
      pull based, GET /config.json, our system could be as simple as that.
4. Support for "factory reset" - what this means I don't know, maybe it's back to a blank box, maybe it undos any user
   applied changes but keeps tool written config (or backdoor config?).
5. Be able to see past version of config, be able to roll back to past versions of config. Maybe roll forward again?
6. Be able to see what changes there are between two sets of config. These could be between versions, between local and
   remote config, between tool and user config, etc.
    1. Diffs between config should be as semantically driven as possible, I don't want to see huge changes listed just
       because items in a set changed order, or whitespace changes.
7. Be able to apply the least amount of change to adjust what I want
    - If I'm changing the PIR timeout for a single space I shouldn't need to restart the node, or even touch other
      spaces.
    - Ideally things would just keep running, but localising change is a good start.
8. Be able to prepare config changes ahead of time/offline, save drafts, etc.
9. Be able to roll out changes across a site (or multiple sites) with confidence
    1. A/B testing maybe
    2. Sequential rollout
    3. Automated testing framework - at the very least be able to see if it worked visually
    4. Automated rollback of config if it didn't work!
10. Be able to update system level config - port numbers, node name, etc.

## Things we've tried before

### sc-bos v1

The config in sc-bos is file based with an "include" concept for loading additional files. Config on disk is loaded into
memory on boot and there are APIs for updating _service_ (drivers, autos, etc) config which updates the config for that
service in memory. Optionally the updated service config can be written to disk in a special location, the original disk
config can be setup to include that user written config location.

The include process is strictly ordered with a first come first served approach. This means any user config dirs must
appear before any non-user config locations in the includes directive of the config file for those user changes to have
an effect.

There is limited central user config management via the edge gateway as the gateway can talk to the individual nodes
APIs and ask them to update their config.

There is no history supported, though we do separate user and tool config files there's no way to see those versions as
any unused config is discarded in memory during include processing.

There is no mechanism for merging user and tool changes within a specific service, this can cause issues as larger
services like the bacnet or tc3dali driver services are monolithic and are configured as a single unit. User changes,
say to update a single lights dali address or even adjust the metadata for a light, will cause the entire services tool
generated config to be ignored from that point on in favour of the user config.

There is no mechanism for changing non service config (ports, node name, CORS settings, etc).

### sc-bos v0.5

The last version of sc-bos (before it was called sc-bos) stored config in a database and distributed it to nodes via a
pull api and a push "there's an update" mechanism. The config was still file based, but the files were stored as blobs
in database rows along with ETag style versioning information.

Tools would generate these config blobs and use the API to write data to the DB. The server would notify nodes that the
config had changed and those nodes would then fetch the new version and apply it locally to update their own systems.
Updating config would typically involve completely reinitialising the node, there was no smaller unit of config update.

Again this project never got around to implementing user driven config changes, all config was written by tools (that
were written by hand). This ended up being quite a pain and ultimately got changed to standard file based config that
was effectively scp'd to the servers and loaded on boot. Now a restart of the node process is needed to apply config
updates.

### Cloud edge gateways (Kahu gateway)

Moving outside of direct sc lineage, we also have cloud edge gateways, appliances that live inside a building with the
goal of connecting it to a cloud service. The config for this system is split into different categories.

Static config, config that is unlikely to change over time, lives on disk in json files and covers config like hosted
web server port number, or which folders to host http static content from. This config can only be edited via these
files and the process needs restarting for changes to be applied. These config files are typically deployed using tools
like Ansible and the source of truth stored in a git repo.

Dynamic config describes the processes domain, which desks are configured, what are their addresses, bus ids, etc. This
config is stored in the cloud and pulled onto the device using cloud apis. The config is editable via cloud hosted admin
pages which sync the config in 'real time' back to the gateway. The gateway is capable of applying these changes without
a restart. There is no mechanism for local editing of this config.

Credential config is stored locally in files on the gateway and can be edited via the embedded admin pages hosted by the
gateway. These credentials can be updated while the gateway is running. The config on disk on the device is the only
location where this config is stored, if you lose this information then it is lost. Only credentials are stored here so
losing them means regenerating new credentials and applying them to the gateway.

A consequence of this type of application is that a cloud connection is required for the appliance to function, so you
may as well use that connection for as much config as you can. A huge benefit is that this allows you to put config into
one and only one location, static config is always on disk, dynamic config is always in the cloud. This bypasses any
conflict resolution issues and simplified the process.

## Bag of ideas

Supporting both tool and manual config generation could be pushed onto the tool to support. While we also want to
support offline config generation we could implement a tool that uses a local cache of config data as a basis for the
modifications it's making. I think ultimately the person who is trying to make the changes needs to see all the
information to make the decision; the decision might be "overwrite all other changes", but until they know what
consequences that has they can't make that decision.

---

Maybe user config should be stored as a diff against the "factory" version of the config. User config that says "set
PIR-02 dali address to 12" should be able to apply to tool config updates in the future. We might have to think about
conflicts but it might work. I think this would also be way more complicated to implement as we'd need to model those
updates and we'd need to introduce some kind of key-framing concept so we don't accidentally implement a change log
database.

---

One option is to use Git for managing and distributing configuration. I think this could come with a bunch of problems
solved for free, however it will also introduce some headaches. Questions that'd need answering if we used Git:

1. Where would the repo live?
2. How would auth be handled?
3. How do we split config across a site? One repo per box, a folder per box, etc.
4. How would things work with multiple avenues for updating the config: user using ops ui changes AC PIR timeout - tool
   pushes new config - ops user edits all PIR timeouts.
5. I don't really want to expose git branching/merging/rebasing/conflicts to users
6. There's no concept of "AC1 has fetched the config" from the pov of the repo