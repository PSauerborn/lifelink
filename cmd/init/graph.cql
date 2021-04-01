CREATE CONSTRAINT unique_module ON (n:Module) ASSERT n.module_name IS UNIQUE;
CREATE CONSTRAINT unique_uid ON (n:User) ASSERT n.uid IS UNIQUE;
CREATE CONSTRAINT unique_email ON (n:User) ASSERT n.email IS UNIQUE;
CREATE CONSTRAINT unique_habit ON (n:Habit) ASSERT n.habit_id IS UNIQUE;
CREATE (u:User {uid: 'lifelink_idp', admin: true, email: 'lifelink@project-gateway.app', created: datetime()});