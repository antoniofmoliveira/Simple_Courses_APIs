package service

import (
	"context"

	"github.com/antoniofmoliveira/courses/db/database"
	"github.com/antoniofmoliveira/courses/dto"

	// "github.com/antoniofmoliveira/courses/dto"
	"github.com/antoniofmoliveira/courses/grpcproto/pb"
)

type CourseService struct {
	pb.UnimplementedCourseServiceServer
	CourseDB database.CourseRepositoryInterface
}

func NewCourseService(courseDB database.CourseRepositoryInterface) *CourseService {
	return &CourseService{
		CourseDB: courseDB,
	}
}

func (c *CourseService) CreateCourse(ctx context.Context, in *pb.CreateCourseRequest) (*pb.Course, error) {
	dtoCourseInputDto := dto.CourseInputDto{Name: in.Name, Description: in.Description, CategoryID: in.CategoryId}
	course, err := c.CourseDB.Create(dtoCourseInputDto)
	if err != nil {
		return nil, err
	}
	return &pb.Course{Id: course.ID, Name: course.Name, Description: course.Description, CategoryId: course.CategoryID}, nil
}

func (c *CourseService) ListCourses(ctx context.Context, in *pb.Blank) (*pb.Courses, error) {
	courses, err := c.CourseDB.FindAll()
	if err != nil {
		return nil, err
	}
	pbCourses := []*pb.Course{}
	for _, course := range courses.Courses {
		pbCourses = append(pbCourses, &pb.Course{Id: course.ID, Name: course.Name, Description: course.Description, CategoryId: course.CategoryID})
	}
	return &pb.Courses{Courses: pbCourses}, nil
}

func (c *CourseService) ListCoursesFromCategory(ctx context.Context, in *pb.ListCoursesFromCategoryRequest) (*pb.Courses, error) {
	courses, err := c.CourseDB.FindByCategoryID(in.CategoryId)
	if err != nil {
		return nil, err
	}
	pbCourses := []*pb.Course{}
	for _, course := range courses.Courses {
		pbCourses = append(pbCourses, &pb.Course{Id: course.ID, Name: course.Name, Description: course.Description, CategoryId: course.CategoryID})
	}
	return &pb.Courses{Courses: pbCourses}, nil
}

func (c *CourseService) GetCourse(ctx context.Context, in *pb.CourseGetRequest) (*pb.Course, error) {
	course, err := c.CourseDB.Find(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Course{Id: course.ID, Name: course.Name, Description: course.Description, CategoryId: course.CategoryID}, nil
}

func (c *CourseService) UpdateCourse(ctx context.Context, in *pb.CourseUpdateRequest) (*pb.Response, error) {
	course := dto.CourseInputDto{ID: in.Id, Name: in.Name, Description: in.Description, CategoryID: in.CategoryId}
	err := c.CourseDB.Update(course)
	if err != nil {
		return nil, err
	}
	return &pb.Response{IsSuccess: true, Message: "Course updated successfully"}, nil
}

func (c *CourseService) DeleteCourse(ctx context.Context, in *pb.CourseDeleteRequest) (*pb.Response, error) {
	err := c.CourseDB.Delete(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Response{IsSuccess: true, Message: "Course deleted successfully"}, nil
}
