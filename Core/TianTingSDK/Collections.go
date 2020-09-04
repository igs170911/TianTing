package TianTingSDK

import "go.mongodb.org/mongo-driver/mongo"

const ColNameMember string = "Member"
const ColNameLUser string = "User"
const ColNameGameUser string = "GameUser"
const ColNameMemberAutoIncrement string = "MemberAutoIncrement"

func (server *TianTingServer) GetMemberCol() *mongo.Collection {
	MongoClient := server.GetDocDb().GetClient()
	return MongoClient.Database(server.CodeName).Collection(ColNameMember)
}

func (server *TianTingServer) GetUserCol() *mongo.Collection {
	MongoClient := server.GetDocDb().GetClient()
	return MongoClient.Database(server.CodeName).Collection(ColNameLUser)
}

func (server *TianTingServer) GetGameUserCol() *mongo.Collection {
	MongoClient := server.GetDocDb().GetClient()
	return MongoClient.Database(server.CodeName).Collection(ColNameGameUser)
}

func (server *TianTingServer) GetMemberAutoIncrementCol() *mongo.Collection {
	MongoClient := server.GetDocDb().GetClient()
	return MongoClient.Database(server.CodeName).Collection(ColNameMemberAutoIncrement)
}