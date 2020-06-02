# UPSERT
`ON CONFLICT ON CONSTRAINT * DO UPDATE`

## Features
- extends the `crudable` interface with an `Upsert` method
- `beforeUpsert` & `afterUpsert` callbacks
  - can't reuse Create/Update callbacks because we don't know which action will be chosen ahead of time
- handles `CreatedAt` & `UpdatedAt` correctly

## Limitations
- only works for postgres (ErrNotImplemented for every other database)
- only works for `int` & `int64` primary keys
  - suggestion: use an `int` as primary key and a `UUID` for public-facing resources (db-internal referencing de-coupled from the public)
- only supports one `UNIQUE CONSTRAINT`
  - requires caller to specify said `CONSTRAINT` for every invocation
  - theoretically, one could use a different constraint for every call to `Upsert` - just not multiple ones
- does not process `associations` (and, by extension, no `eager` loading)
- caller has to explicitly specify if `PrimaryKey` conflicts should be considered
> despite these severe limitations, it has served me well so far ;)

## Solutions (?)
- execute Upsert for all `associations`
- automatically use the `PrimaryKey` as `CONSTRAINT` if
  - it's actually defined (!= 0)
  - OR no UniqueKey was submitted

## Add to a gobuffalo project
Append the following to the end of `go.mod`:\
`replace github.com/gobuffalo/pop/v5 => github.com/julius-b/pop/v5 v5.1.3-upsert`\
v4: `replace github.com/gobuffalo/pop => github.com/julius-b/pop v4.13.1-upsert+incompatible`

## Usage
`fizz` can only create unique **indicies**, so you have to define an additional unique **constraint** in your `models.up.fizz` file:
```golang
// fizz default
add_index("contacts", ["contact_id", "installation_id"], {"unique": true})

// new UNIQUE CONSTRAINT with same columns 
sql("ALTER TABLE contacts ADD UNIQUE(contact_id, installation_id)")
```
> fizz doesn't actually support comments so you'll have to remove them...

then, store UniqueKey's name for every model, e.g:
```golang
const ContactUniqueKey = "contacts_contact_id_installation_id_key"
```

finally, the usage is *almost* the same as with Create / Update:
```golang
// the last bool defines if a conflict with the PrimaryKey should trigger an Upsert
//   - only considered if it's value is not 0
// If you have defined a special constraint then it should be false, otherwise probably true
verrs, err := tx.ValidateAndUpsert(contact, models.ContactUniqueKey, false)
if err != nil {
  c.Logger().Warnf("contacts/Create - upsert failed: %v", err)
  return err
}

// if verrs.HasAny() ...
```
