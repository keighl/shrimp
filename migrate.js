r.dbCreate('shrimp');

r.db('shrimp').tableCreate('todos');
r.db('shrimp').table('todos').indexCreate('user_id');

r.db('shrimp').tableCreate('password_resets');
r.db('shrimp').table('password_resets').indexCreate('user_id');

r.db('shrimp').tableCreate('users');
r.db('shrimp').table('users').indexCreate('api_token');
r.db('shrimp').table('users').indexCreate('email');
