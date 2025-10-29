// services 包包含所有业务逻辑的实现
// 这个文件专门处理用户相关的业务逻辑
package services

// 导入所需的包
import (
	"blogV2/internal/models"     // 数据模型包，包含数据库实体定义
	"blogV2/pkg/utils"          // 工具包，包含JWT生成等工具函数
	"errors"                     // 错误处理包，用于创建自定义错误

	"golang.org/x/crypto/bcrypt"    // 密码加密包，用于密码的bcrypt加密
	"gorm.io/gorm"                  // GORM ORM框架，用于数据库操作
)

// UserService 用户服务结构体
// 负责处理所有与用户相关的业务逻辑
// 遵循服务层模式，作为控制层和模型层之间的桥梁
type UserService struct {
	db *gorm.DB // 数据库连接实例，用于执行数据库操作
}

// NewUserService 创建用户服务实例
// 这是服务的工厂函数，用于创建并初始化UserService实例
// 参数:
//   db - 数据库连接实例
// 返回值:
//   *UserService - 初始化好的用户服务实例
func NewUserService(db *gorm.DB) *UserService {
	// 返回新创建的UserService实例
	// 将传入的数据库连接赋值给服务的db字段
	return &UserService{db: db}
}

// CreateUser 创建新用户
// 这个方法处理用户注册的业务逻辑，包括数据验证和密码加密
// 参数:
//   user - 用户模型实例，包含用户注册信息
// 返回值:
//   error - 如果创建过程中发生错误则返回错误信息，否则返回nil
func (s *UserService) CreateUser(user *models.User) error {
	// 检查用户名是否已存在
	// 查询数据库中是否已存在相同用户名的用户
	var existingUser models.User
	result := s.db.Where("username = ?", user.Username).First(&existingUser)
	
	// 如果查询成功（找到了记录），说明用户名已存在
	if result.Error == nil {
		// 返回自定义错误，告知客户端用户名已被占用
		return errors.New("用户名已存在")
	}
	
	// 检查邮箱是否已存在
	// 查询数据库中是否已存在相同邮箱的用户
	result = s.db.Where("email = ?", user.Email).First(&existingUser)
	
	// 如果查询成功（找到了记录），说明邮箱已被占用
	if result.Error == nil {
		// 返回自定义错误，告知客户端邮箱已被注册
		return errors.New("邮箱已被注册")
	}

	// 创建用户记录
	// 使用GORM的Create方法将用户数据插入数据库
	// 注意：密码会在模型的BeforeSave钩子中自动加密
	if err := s.db.Create(user).Error; err != nil {
		// 如果创建过程中发生错误，返回错误信息
		return err
	}

	// 用户创建成功，返回nil表示没有错误
	return nil
}

// Login 用户登录验证
// 这个方法处理用户登录的业务逻辑，包括密码验证和JWT token生成
// 参数:
//   username - 用户名
//   password - 密码（明文）
// 返回值:
//   *models.User - 登录成功的用户信息
//   string - 生成的JWT token
//   error - 如果登录过程中发生错误则返回错误信息
func (s *UserService) Login(username, password string) (*models.User, string, error) {
	// 根据用户名查询用户信息
	var user models.User
	
	// 使用GORM的Where方法根据用户名查询用户
	result := s.db.Where("username = ?", username).First(&user)
	
	// 如果查询失败（用户不存在）
	if result.Error != nil {
		// 返回自定义错误，告知客户端用户名或密码错误（不具体说明是用户名错误）
		return nil, "", errors.New("用户名或密码错误")
	}

	// 验证密码是否正确
	// 使用bcrypt包比较输入的密码和数据库中存储的加密密码
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// 如果密码不匹配，返回自定义错误
		return nil, "", errors.New("用户名或密码错误")
	}

	// 生成JWT token
	tokenString, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		// 如果token生成失败，返回错误信息
		return nil, "", err
	}

	// 登录成功，返回用户信息和JWT token
	return &user, tokenString, nil
}

// GetUserByID 根据用户ID获取用户信息
// 这个方法根据用户ID查询用户详细信息
// 参数:
//   id - 用户ID
// 返回值:
//   *models.User - 查询到的用户信息
//   error - 如果查询过程中发生错误则返回错误信息
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	// 声明用户变量，用于存储查询结果
	var user models.User
	
	// 使用GORM的First方法根据主键查询用户
	result := s.db.First(&user, id)
	
	// 如果查询失败，返回错误信息
	if result.Error != nil {
		return nil, result.Error
	}

	// 查询成功，返回用户信息
	return &user, nil
}

// UpdateUser 更新用户信息
// 这个方法处理用户信息更新的业务逻辑
// 参数:
//   user - 包含更新后信息的用户模型实例
// 返回值:
//   error - 如果更新过程中发生错误则返回错误信息
func (s *UserService) UpdateUser(user *models.User) error {
	// 使用GORM的Save方法更新用户记录
	// Save方法会根据主键判断是更新还是创建
	if err := s.db.Save(user).Error; err != nil {
		return err
	}
	
	// 更新成功，返回nil
	return nil
}

// DeleteUser 删除用户
// 这个方法处理用户删除的业务逻辑
// 参数:
//   id - 要删除的用户ID
// 返回值:
//   error - 如果删除过程中发生错误则返回错误信息
func (s *UserService) DeleteUser(id uint) error {
	// 使用GORM的Delete方法删除用户记录
	if err := s.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	
	// 删除成功，返回nil
	return nil
}