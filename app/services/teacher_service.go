package services

import (
	"classroom/app/configs"
	"classroom/app/models"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// interfáz

type TeachersService interface {
	CreateTeacher(ctx context.Context, teacherRequest *models.TeacherCreateRequest) (*models.Teacher, error)
	UpdateTeacher(ctx context.Context, teacherRequest *models.TeacherCreateRequest, ID primitive.ObjectID) (*models.Teacher, error)
}

type teachersServiceImpl struct{}

func NewTeachersService() TeachersService {
	return &teachersServiceImpl{}
}

func ValidateTeacherData(teacherRequest *models.TeacherCreateRequest) (*models.Teacher, error) {
	// Validaciones básicas de campos requeridos
	if teacherRequest.Names == "" {
		return nil, errors.New("el nombre es requerido")
	}
	if teacherRequest.LastNames == "" {
		return nil, errors.New("los apellidos son requeridos")
	}
	if teacherRequest.DocumentNumber == "" {
		return nil, errors.New("el número de documento es requerido")
	}
	if teacherRequest.Code == "" {
		return nil, errors.New("el código de profesor es requerido")
	}

	// Validar que el documentTypeID sea un ObjectID válido
	if teacherRequest.DocumentTypeID == "" {
		return nil, errors.New("el tipo de documento es requerido")
	}
	if !primitive.IsValidObjectID(teacherRequest.DocumentTypeID) {
		return nil, errors.New("el ID de tipo de documento no es válido")
	}

	// Convertir documentTypeID a ObjectID
	docTypeID, err := primitive.ObjectIDFromHex(teacherRequest.DocumentTypeID)
	if err != nil {
		return nil, fmt.Errorf("error convirtiendo documentTypeID: %v", err)
	}

	// Crear estructura Teacher
	teacher := &models.Teacher{
		Person: models.Person{
			Type:           "teacher",
			Names:          teacherRequest.Names,
			LastNames:      teacherRequest.LastNames,
			DocumentNumber: teacherRequest.DocumentNumber,
			DocumentTypeID: docTypeID,
			ImageURL:       teacherRequest.ImageURL,
			Created:        time.Now(),
			Updated:        time.Now(),
		},
		Code: teacherRequest.Code,
	}

	return teacher, nil
}

func (s *teachersServiceImpl) CreateTeacher(ctx context.Context, teacherRequest *models.TeacherCreateRequest) (*models.Teacher, error) {
	// Validar y crear estructura Teacher
	teacher, err := ValidateTeacherData(teacherRequest)
	if err != nil {
		return nil, err
	}

	// Insertar en la colección de persons (ya que Teacher hereda de Person)
	collection := configs.DB.Collection("persons")
	result, err := collection.InsertOne(ctx, teacher)
	if err != nil {
		return nil, fmt.Errorf("error al crear profesor: %v", err)
	}

	// Asignar el ID generado
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		teacher.ID = oid
		return teacher, nil
	}

	return nil, errors.New("no se pudo obtener el ID generado")
}

func (s *teachersServiceImpl) UpdateTeacher(ctx context.Context, teacherRequest *models.TeacherCreateRequest, ID primitive.ObjectID) (*models.Teacher, error) {
	// 1. Validar los datos de entrada
	teacher, err := ValidateTeacherData(teacherRequest)
	if err != nil {
		return nil, err
	}

	// 2. Preparar la actualización
	collection := configs.DB.Collection("persons")
	update := bson.M{
		"$set": bson.M{
			"type":             teacher.Type,
			"names":            teacher.Names,
			"last_names":       teacher.LastNames,
			"document_number":  teacher.DocumentNumber,
			"document_type_id": teacher.DocumentTypeID,
			"image_url":        teacher.ImageURL,
			"code":             teacher.Code,
			"updated":          time.Now(),
		},
	}

	// 3. Ejecutar la actualización
	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": ID, "type": "teacher"}, // Filtro por ID y tipo teacher
		update,
	)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar profesor: %v", err)
	}

	// 4. Verificar si se actualizó algún documento
	if result.MatchedCount == 0 {
		return nil, errors.New("profesor no encontrado")
	}

	// 5. Asignar el ID original y devolver el profesor actualizado
	teacher.ID = ID
	return teacher, nil
}
