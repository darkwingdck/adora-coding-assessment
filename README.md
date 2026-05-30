# coding-assessment

# draft notes
- Trade off sqlite/postgresql - triggers, enums, separate container
- It's weird to have text as user_id, should consider to make it uuid. If text is id from another service - this is also bad. Also, it's weird for us to call carrier route with out db's user id. It's better to have two of them - of is just db's id (uuid), and another is whatever we want. For example, mobile_carrier_id
- Would be nice to add golang-migrate, but that's not a first priority rn
- Add gomock for tests (?)