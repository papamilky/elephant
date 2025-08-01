# Elephant - cuz it's phat.

`elephant` is a service providing data and actions via various data-providers. It is meant to be a backend to create f.e. custom launchers.

## Current State

The project just started and is therefore highly wip.

## Current Providers

- `desktopapplications`
- `files`

## Quick-Guide

1. You need: `elephant`
2. ... a provider
3. something to make unix socket calls with

```
mkdir ~/.config/elephant
mkdir ~/.config/elephant/providers
git clone https://github.com/abenz1267/elephant && cd elephant/cmd && go install elephant.go
cd ../internals/providers/desktopapplications
go build -buildmode=plugin && cp desktopapplications.so ~/.config/elephant/providers/
```

Once you have this setup, you can start using `elephant`.

### How to query:

`qid = queryid, iid = iterationid`

1. Open socket connection with f.e. `nc -U /tmp/elephant.sock`
2. You can now query with: `query;desktopapplications;firefox`
3. You'll retrieve a `qid;<qid>;<iid>`
4. You'll retrieve a list of entries:

```
1;1;desktopapplications;/usr/share/applications/firefox-developer-edition.desktop;Firefox Developer Edition;Web Browser;firefox-developer-edition;6,5,4,3,2,1,0;0;text
```

To break this down:

```
QID;IID;PROVIDER;IDENTIFIER;TEXT;SUBTEXT;ICON;POSITIONS OF FUZZY MATCH;STARTING POSITION OF FUZZY MATCH;FUZZY MATCH FIELDNAME
```

### How to subscribe to changes:

Subscribe to a query with an interval:

```
subscribe;1000;desktopapplications;foot
```

Subscribe to realtime data-changes of a provider:

```
subscribe;0;desktopapplications;
```

You'll retrieve `sid;100000001;changed` as status updates.

You can unsubscribe again by doing `unsubscribe;100000001`.

### How to activate:

```
activate;1;desktopapplications;firefox-developer-edition.desktop;
```

To break this down:

```
COMMAND;QID;PROVIDER;IDENTIFIER;ACTION
```

If there's no action, you'll still need a trailing `;`.
