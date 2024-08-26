import redis
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from sqlalchemy.sql import text  # 导入 text

# 创建数据库引擎
# 创建 Redis 连接
r = redis.Redis(host='localhost', port=6379, db=0)

SQLALCHEMY_DATABASE_URL = "mysql+pymysql://caller:imXAGcrGSdaxH8RT@localhost/caller"
# SQLALCHEMY_DATABASE_URL =  "mysql+pymysql://root:root@localhost/call"
engine = create_engine(SQLALCHEMY_DATABASE_URL,pool_pre_ping=True,pool_recycle=1800,echo=False)
# engine = create_engine('your_database_connection_string')

# 创建会话
Session = sessionmaker(bind=engine)
session = Session()

# 执行 MySQL 语句，使用 text 函数
results = session.execute(text("SELECT group_client.group_name,group_client.client_id from group_client"))

# 连接 Redis
r = redis.Redis(host='localhost', port=6379)

# 获取所有以 "group_name_" 为前缀的键
keys = r.keys('group_name_*')

# 遍历这些键，并判断其类型是否为列表，如果是则删除
for key in keys:
 
    r.delete(key)
    print("删除",key)

# 处理结果并存入 Redis
for row in results:
    group_name = row[0]
    client_id = row[1]
    key = "group_name_" + group_name
    print(key,client_id)
    r.rpush(key, client_id)

# 关闭会话和 Redis 连接
session.close()
r.close()