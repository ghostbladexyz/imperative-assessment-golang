# Key progress by exercise identity

Exercise Keys are immutable, source-qualified strings, while numeric identifiers remain changeable Curriculum Positions. Saved Exercise Progress and pass receipts use Exercise Keys so curriculum reordering does not rename work. The checkpoint catalogue is a replacement for the former curriculum, so progress schema 6 intentionally does not migrate schema-4 or schema-5 records; incompatible saved progress starts clean rather than applying an obsolete positional mapping.
