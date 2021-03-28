CREATE CONSTRAINT unique_module ON (n:Module) ASSERT n.module_name IS UNIQUE;
CREATE CONSTRAINT unique_uid ON (n:User) ASSERT n.uid IS UNIQUE;
CREATE CONSTRAINT unique_email ON (n:User) ASSERT n.email IS UNIQUE;