package storage

type RedisStorage struct {
	dbs       *DatabaseService
	keyPrefix string
}

func (rs *RedisStorage) SetDatabaseService(s *DatabaseService) {
	rs.dbs = s
}

func (rs *RedisStorage) SetKeyPrefix(kp string) {
	rs.keyPrefix = kp
}

func (rs *RedisStorage) GetDatabaseService() *DatabaseService {

	return rs.dbs
}

func (rs *RedisStorage) GetKeyPrefix() string {
	return rs.keyPrefix
}
