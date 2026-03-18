INSERT INTO words (word, date) VALUES
  ('hello', CURRENT_DATE),
  ('world', CURRENT_DATE + 1),
  ('svelte', CURRENT_DATE + 2),
  ('router', CURRENT_DATE + 3),
  ('query', CURRENT_DATE + 4)
ON CONFLICT DO NOTHING;
