package main

import (
	"bytes"
	"context"
	db "jiaojiao/database"
	content "jiaojiao/srv/content/proto"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestContentCreate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := db.MongoDatabase.Collection("content")

	var req content.ContentCreateRequest

	tf := func(status content.ContentCreateResponse_Status, success bool) (string, string) {
		var s srv
		var rsp content.ContentCreateResponse
		So(s.Create(context.TODO(), &req, &rsp), ShouldBeNil)
		So(rsp.Status, ShouldEqual, status)
		if success {
			So(rsp.ContentID, ShouldNotBeBlank)
			So(rsp.ContentToken, ShouldNotBeBlank)
			So(rsp.FileID, ShouldNotBeBlank)
		} else {
			So(rsp.ContentID, ShouldBeBlank)
			So(rsp.ContentToken, ShouldBeBlank)
			So(rsp.FileID, ShouldBeBlank)
		}
		return rsp.ContentID, rsp.ContentToken
	}

	Convey("Test content create", t, func() {
		req.ContentID = ""
		req.ContentToken = ""
		req.Content = []byte("valid_file")
		req.Type = content.Type_VIDEO
		id, token := tf(content.ContentCreateResponse_SUCCESS, true)

		req.ContentID = id
		req.ContentToken = token
		req.Content = []byte("valid_file")
		req.Type = content.Type_PICTURE
		tf(content.ContentCreateResponse_SUCCESS, true)
		defer func() {
			var res result
			cid, err := primitive.ObjectIDFromHex(id)
			So(err, ShouldBeNil)
			err = collection.FindOneAndDelete(ctx, bson.D{
				{"_id", cid},
				{"token", token},
			}).Decode(&res)
			So(err, ShouldBeNil)
		}()

		req.ContentID = id
		req.ContentToken = "12463-25897fsfs-5232"
		tf(content.ContentCreateResponse_INVALID_TOKEN, false)

		req.ContentID = "1234"
		req.ContentToken = token
		tf(content.ContentCreateResponse_INVALID_TOKEN, false)

		req.ContentID = ""
		req.ContentToken = token
		tf(content.ContentCreateResponse_INVALID_PARAM, false)

		req.ContentID = id
		req.ContentToken = ""
		tf(content.ContentCreateResponse_INVALID_PARAM, false)

		req.ContentID = id
		req.ContentToken = token
		req.Type = 0
		tf(content.ContentCreateResponse_INVALID_PARAM, false)

		req.Type = content.Type_PICTURE
		req.Content = []byte{0}
		tf(content.ContentCreateResponse_INVALID_PARAM, false)

	})
}

func TestCreateTag(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := db.MongoDatabase.Collection("content")

	var req content.ContentCreateTagRequest

	tf := func(status content.ContentCreateTagResponse_Status) (string, string) {
		var s srv
		var rsp content.ContentCreateTagResponse
		So(s.CreateTag(context.TODO(), &req, &rsp), ShouldBeNil)
		So(rsp.Status, ShouldEqual, status)
		if rsp.Status == content.ContentCreateTagResponse_SUCCESS {
			So(rsp.ContentID, ShouldNotBeBlank)
			So(rsp.ContentToken, ShouldNotBeBlank)

			cid, err := primitive.ObjectIDFromHex(rsp.ContentID)
			So(err, ShouldBeNil)
			var res result
			err = collection.FindOne(ctx, bson.D{
				{"_id", cid},
			}).Decode(&res)
			So(err, ShouldBeNil)
			So(res.Tags, ShouldResemble, req.Tags)
		} else {
			So(rsp.ContentID, ShouldBeBlank)
			So(rsp.ContentToken, ShouldBeBlank)
		}
		return rsp.ContentID, rsp.ContentToken
	}

	Convey("Test content create tag", t, func() {
		req.ContentID = ""
		req.ContentToken = ""
		req.Tags = []string{"book", "cook"}
		id, token := tf(content.ContentCreateTagResponse_SUCCESS)
		defer func() {
			var res result
			cid, err := primitive.ObjectIDFromHex(id)
			So(err, ShouldBeNil)
			err = collection.FindOneAndDelete(ctx, bson.D{
				{"_id", cid},
				{"token", token},
			}).Decode(&res)
			So(err, ShouldBeNil)
		}()

		req.ContentID = id
		req.ContentToken = token
		req.Tags = []string{"dook"}
		tf(content.ContentCreateTagResponse_SUCCESS)

		req.ContentToken = "456"
		tf(content.ContentCreateTagResponse_INVALID_TOKEN)

		req.ContentID = "123"
		req.ContentToken = token
		tf(content.ContentCreateTagResponse_INVALID_TOKEN)

		req.ContentID = ""
		req.ContentToken = token
		tf(content.ContentCreateTagResponse_INVALID_PARAM)

		req.ContentID = id
		req.ContentToken = ""
		tf(content.ContentCreateTagResponse_INVALID_PARAM)

		req.ContentID = id
		req.ContentToken = token
		req.Tags = []string{}
		tf(content.ContentCreateTagResponse_INVALID_PARAM)
	})
}

func TestContentDelete(t *testing.T) {
	tf := func(status content.ContentDeleteResponse_Status, id string, token string, fid string) {
		var s srv
		var rsp content.ContentDeleteResponse
		So(s.Delete(context.TODO(), &content.ContentDeleteRequest{
			ContentID:    id,
			ContentToken: token,
			FileID:       fid,
		}, &rsp), ShouldBeNil)
		So(rsp.Status, ShouldEqual, status)
	}

	Convey("Test content delete", t, func() {
		bucket, err := gridfs.NewBucket(db.MongoDatabase)
		So(err, ShouldBeNil)
		objID, err := bucket.UploadFromStream("", bytes.NewReader([]byte("valid")))
		So(err, ShouldBeNil)
		fid := objID.Hex()
		defer func() { So(bucket.Delete(objID), ShouldNotBeBlank) }()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		collection := db.MongoDatabase.Collection("content")
		token := uuid.NewV4().String()
		res, err := collection.InsertOne(ctx, bson.M{
			"token": token,
			"files": bson.A{
				bson.M{
					"fileID": fid,
					"type":   1,
				}},
		})
		So(err, ShouldBeNil)
		id := res.InsertedID.(primitive.ObjectID).Hex()
		defer func() {
			err = collection.FindOneAndDelete(ctx, bson.D{
				{"_id", res.InsertedID.(primitive.ObjectID)},
				{"token", token},
			}).Decode(&res)
			So(err, ShouldNotBeBlank)
		}()

		tf(content.ContentDeleteResponse_INVALID_TOKEN, "123123", token, fid)
		tf(content.ContentDeleteResponse_INVALID_TOKEN, id, "123123", fid)
		tf(content.ContentDeleteResponse_NOT_FOUND, id, token, "123123")
		tf(content.ContentDeleteResponse_INVALID_PARAM, "", token, fid)
		tf(content.ContentDeleteResponse_INVALID_PARAM, id, "", fid)
		tf(content.ContentDeleteResponse_INVALID_PARAM, id, token, "")
		tf(content.ContentDeleteResponse_SUCCESS, id, token, fid)
		tf(content.ContentDeleteResponse_SUCCESS, id, token, "")
	})
}

func TestContentQuery(t *testing.T) {
	tf := func(status content.ContentQueryResponse_Status, id string, success bool) {
		var s srv
		var rsp content.ContentQueryResponse
		So(s.Query(context.TODO(), &content.ContentQueryRequest{
			ContentID: id,
		}, &rsp), ShouldBeNil)
		So(rsp.Status, ShouldEqual, status)
		if success {
			So(rsp.ContentToken, ShouldNotBeBlank)
			So(rsp.Files, ShouldNotBeBlank)
			So(rsp.Tags, ShouldNotBeBlank)
		} else {
			So(rsp.ContentToken, ShouldBeBlank)
			So(rsp.Files, ShouldBeBlank)
			So(rsp.Tags, ShouldBeBlank)
		}
	}

	Convey("Test content query", t, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		collection := db.MongoDatabase.Collection("content")
		token := uuid.NewV4().String()
		res, err := collection.InsertOne(ctx, bson.M{
			"token": token,
			"files": bson.A{
				bson.M{
					"fileID": "012345678901234567891234",
					"type":   1,
				}},
		})
		So(err, ShouldBeNil)
		id := res.InsertedID.(primitive.ObjectID).Hex()
		defer func() {
			err = collection.FindOneAndDelete(ctx, bson.D{
				{"_id", res.InsertedID.(primitive.ObjectID)},
				{"token", token},
			}).Decode(&res)
			So(err, ShouldBeNil)
		}()

		tf(content.ContentQueryResponse_SUCCESS, id, true)
		tf(content.ContentQueryResponse_NOT_FOUND, "123123", false)
		tf(content.ContentQueryResponse_INVALID_PARAM, "", false)
	})
}

func TestContentUpdate(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := db.MongoDatabase.Collection("content")

	var req content.ContentUpdateRequest

	tf := func(status content.ContentUpdateResponse_Status, success bool) string {
		var s srv
		var rsp content.ContentUpdateResponse
		So(s.Update(context.TODO(), &req, &rsp), ShouldBeNil)
		So(rsp.Status, ShouldEqual, status)
		if success {
			So(rsp.FileID, ShouldNotBeBlank)
		} else {
			So(rsp.FileID, ShouldBeBlank)
		}
		return rsp.FileID
	}

	Convey("Test content update", t, func() {
		bucket, err := gridfs.NewBucket(db.MongoDatabase)
		So(err, ShouldBeNil)
		objID, err := bucket.UploadFromStream("", bytes.NewReader([]byte("valid")))
		So(err, ShouldBeNil)
		fid := objID.Hex()
		defer func() { So(bucket.Delete(objID), ShouldBeNil) }()

		token := uuid.NewV4().String()
		res, err := collection.InsertOne(ctx, bson.M{
			"token": token,
			"files": bson.A{
				bson.M{
					"fileID": fid,
					"type":   1,
				}},
		})
		So(err, ShouldBeNil)
		id := res.InsertedID.(primitive.ObjectID).Hex()
		defer func() {
			err = collection.FindOneAndDelete(ctx, bson.D{
				{"_id", res.InsertedID.(primitive.ObjectID)},
				{"token", token},
			}).Decode(&res)
			So(err, ShouldBeNil)
		}()

		req.ContentID = id
		req.ContentToken = token
		req.FileID = fid
		req.Content = []byte("valid_file")
		req.Type = content.Type_PICTURE
		req.FileID = tf(content.ContentUpdateResponse_SUCCESS, true)

		req.Type = content.Type_VIDEO
		fid = tf(content.ContentUpdateResponse_SUCCESS, true)

		tf(content.ContentUpdateResponse_NOT_FOUND, false)

		req.FileID = fid
		req.ContentID = id
		req.ContentToken = "12463-25897fsfs-5232"
		tf(content.ContentUpdateResponse_INVALID_TOKEN, false)

		req.ContentID = "1234"
		req.ContentToken = token
		tf(content.ContentUpdateResponse_INVALID_TOKEN, false)

		req.ContentID = ""
		req.ContentToken = token
		tf(content.ContentUpdateResponse_INVALID_PARAM, false)

		req.ContentID = id
		req.ContentToken = ""
		tf(content.ContentUpdateResponse_INVALID_PARAM, false)

		req.ContentID = id
		req.ContentToken = token
		req.FileID = ""
		tf(content.ContentUpdateResponse_INVALID_PARAM, false)

		req.Type = content.Type_PICTURE
		req.Content = []byte{0}
		tf(content.ContentUpdateResponse_INVALID_PARAM, false)

		req.ContentID = id
		req.ContentToken = token
		req.Type = 0
		tf(content.ContentUpdateResponse_INVALID_TYPE, false)

		req.Type = 3
		tf(content.ContentUpdateResponse_INVALID_TYPE, false)
	})
}

func TestContentCheck(t *testing.T) {
	tf := func(status content.ContentCheckResponse_Status, id string, token string) {
		var s srv
		var rsp content.ContentCheckResponse
		So(s.Check(context.TODO(), &content.ContentCheckRequest{
			ContentID:    id,
			ContentToken: token,
		}, &rsp), ShouldBeNil)
		So(rsp.Status, ShouldEqual, status)
	}

	Convey("Test content check", t, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		collection := db.MongoDatabase.Collection("content")
		token := uuid.NewV4().String()
		res, err := collection.InsertOne(ctx, bson.M{
			"token": token,
			"files": bson.A{
				bson.M{
					"fileID": "012345678901234567891234",
					"type":   1,
				}},
		})
		So(err, ShouldBeNil)
		id := res.InsertedID.(primitive.ObjectID).Hex()
		defer func() {
			err = collection.FindOneAndDelete(ctx, bson.D{
				{"_id", res.InsertedID.(primitive.ObjectID)},
				{"token", token},
			}).Decode(&res)
			So(err, ShouldBeNil)
		}()

		tf(content.ContentCheckResponse_VALID, id, token)
		tf(content.ContentCheckResponse_INVALID, id, "123123")
		tf(content.ContentCheckResponse_INVALID, "123123", token)
		tf(content.ContentCheckResponse_INVALID, "123123", "123123")
		tf(content.ContentCheckResponse_INVALID_PARAM, id, "")
		tf(content.ContentCheckResponse_INVALID_PARAM, "", token)
		tf(content.ContentCheckResponse_INVALID_PARAM, "", "123123")
		tf(content.ContentCheckResponse_INVALID_PARAM, "123123", "")
		tf(content.ContentCheckResponse_INVALID_PARAM, "", "")
	})
}

func TestMain(m *testing.M) {
	main()
	m.Run()
}
