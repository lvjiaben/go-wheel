package crypto

import (
	"strings"
	"sync"
	"testing"
)

// TestGenerateSalt 测试盐值生成
func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("生成盐值失败: %v", err)
	}

	if len(salt) != 32 { // 16 bytes = 32 hex chars
		t.Errorf("盐值长度应该是 32，实际: %d", len(salt))
	}

	// 验证是否为有效的十六进制字符串
	for _, c := range salt {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("盐值包含非法字符: %c", c)
		}
	}
}

// TestGenerateSaltUniqueness 测试盐值唯一性
func TestGenerateSaltUniqueness(t *testing.T) {
	salts := make(map[string]bool)
	count := 1000

	for i := 0; i < count; i++ {
		salt, err := GenerateSalt()
		if err != nil {
			t.Fatalf("生成盐值失败: %v", err)
		}
		if salts[salt] {
			t.Errorf("生成了重复的盐值: %s", salt)
		}
		salts[salt] = true
	}

	if len(salts) != count {
		t.Errorf("应该生成 %d 个唯一盐值，实际: %d", count, len(salts))
	}
}

// TestHashPassword 测试密码哈希
func TestHashPassword(t *testing.T) {
	password := "test123456"
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("生成盐值失败: %v", err)
	}

	hash, err := HashPassword(password, salt)
	if err != nil {
		t.Fatalf("哈希密码失败: %v", err)
	}

	if hash == "" {
		t.Error("哈希值不应该为空")
	}

	if hash == password {
		t.Error("哈希值不应该等于原始密码")
	}

	// 验证哈希值格式（bcrypt 哈希以 $2a$ 或 $2b$ 开头）
	if !strings.HasPrefix(hash, "$2a$") && !strings.HasPrefix(hash, "$2b$") {
		t.Errorf("哈希值格式不正确: %s", hash)
	}
}

// TestHashPasswordConsistency 测试相同输入产生相同哈希
func TestHashPasswordConsistency(t *testing.T) {
	password := "test123456"
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("生成盐值失败: %v", err)
	}

	hash1, err1 := HashPassword(password, salt)
	if err1 != nil {
		t.Fatalf("第一次哈希失败: %v", err1)
	}

	hash2, err2 := HashPassword(password, salt)
	if err2 != nil {
		t.Fatalf("第二次哈希失败: %v", err2)
	}

	// bcrypt 每次都会生成不同的哈希（因为内部有随机盐）
	// 但我们可以验证两个哈希都能验证相同的密码
	if !PasswordVerifyWithSalt(password, salt, hash1) {
		t.Error("第一个哈希验证失败")
	}

	if !PasswordVerifyWithSalt(password, salt, hash2) {
		t.Error("第二个哈希验证失败")
	}
}

// TestVerifyPassword 测试密码验证
func TestVerifyPassword(t *testing.T) {
	password := "test123456"
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("生成盐值失败: %v", err)
	}

	hash, err := HashPassword(password, salt)
	if err != nil {
		t.Fatalf("哈希密码失败: %v", err)
	}

	// 测试正确密码
	if !PasswordVerifyWithSalt(password, salt, hash) {
		t.Error("正确密码验证失败")
	}

	// 测试错误密码
	if PasswordVerifyWithSalt("wrongpassword", salt, hash) {
		t.Error("错误密码不应该验证通过")
	}

	// 测试错误盐值
	wrongSalt, _ := GenerateSalt()
	if PasswordVerifyWithSalt(password, wrongSalt, hash) {
		t.Error("错误盐值不应该验证通过")
	}
}

// TestVerifyPasswordEdgeCases 测试边界情况
func TestVerifyPasswordEdgeCases(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Fatalf("生成盐值失败: %v", err)
	}

	// 测试空密码
	emptyHash, err := HashPassword("", salt)
	if err != nil {
		t.Fatalf("空密码哈希失败: %v", err)
	}
	if !PasswordVerifyWithSalt("", salt, emptyHash) {
		t.Error("空密码验证失败")
	}

	// 测试长密码（bcrypt 限制为 72 字节，salt 是 32 字节，所以密码最多 40 字节）
	longPassword := strings.Repeat("a", 40)
	longHash, err := HashPassword(longPassword, salt)
	if err != nil {
		t.Fatalf("长密码哈希失败: %v", err)
	}
	if !PasswordVerifyWithSalt(longPassword, salt, longHash) {
		t.Error("长密码验证失败")
	}

	// 测试特殊字符密码
	specialPassword := "!@#$%^&*()_+-=[]{}|;:',.<>?/~`"
	specialHash, err := HashPassword(specialPassword, salt)
	if err != nil {
		t.Fatalf("特殊字符密码哈希失败: %v", err)
	}
	if !PasswordVerifyWithSalt(specialPassword, salt, specialHash) {
		t.Error("特殊字符密码验证失败")
	}

	// 测试中文密码
	chinesePassword := "测试密码123"
	chineseHash, err := HashPassword(chinesePassword, salt)
	if err != nil {
		t.Fatalf("中文密码哈希失败: %v", err)
	}
	if !PasswordVerifyWithSalt(chinesePassword, salt, chineseHash) {
		t.Error("中文密码验证失败")
	}
}

// TestPasswordConcurrency 测试并发安全
func TestPasswordConcurrency(t *testing.T) {
	var wg sync.WaitGroup
	password := "test123456"
	iterations := 100

	// 并发生成盐值
	salts := make([]string, iterations)
	saltErrors := make([]error, iterations)
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			salts[index], saltErrors[index] = GenerateSalt()
		}(i)
	}
	wg.Wait()

	// 验证所有盐值生成成功
	for i, err := range saltErrors {
		if err != nil {
			t.Fatalf("并发生成盐值 %d 失败: %v", i, err)
		}
	}

	// 验证所有盐值都是唯一的
	saltMap := make(map[string]bool)
	for _, salt := range salts {
		if saltMap[salt] {
			t.Error("并发生成了重复的盐值")
		}
		saltMap[salt] = true
	}

	// 并发哈希密码
	hashes := make([]string, iterations)
	errors := make([]error, iterations)
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			hashes[index], errors[index] = HashPassword(password, salts[index])
		}(i)
	}
	wg.Wait()

	// 验证所有哈希都成功
	for i, err := range errors {
		if err != nil {
			t.Errorf("并发哈希 %d 失败: %v", i, err)
		}
	}

	// 并发验证密码
	results := make([]bool, iterations)
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index] = PasswordVerifyWithSalt(password, salts[index], hashes[index])
		}(i)
	}
	wg.Wait()

	// 验证所有验证都成功
	for i, result := range results {
		if !result {
			t.Errorf("并发验证 %d 失败", i)
		}
	}
}

// TestDifferentSaltsDifferentHashes 测试不同盐值产生不同哈希
func TestDifferentSaltsDifferentHashes(t *testing.T) {
	password := "test123456"
	salt1, _ := GenerateSalt()
	salt2, _ := GenerateSalt()

	hash1, _ := HashPassword(password, salt1)
	hash2, _ := HashPassword(password, salt2)

	// 虽然密码相同，但盐值不同，所以哈希应该不同
	// 注意：这里我们比较的是加盐后的密码，bcrypt 会再次加盐
	// 所以即使盐值相同，bcrypt 的哈希也可能不同
	// 但我们可以验证两个哈希都能用各自的盐值验证
	if !PasswordVerifyWithSalt(password, salt1, hash1) {
		t.Error("hash1 验证失败")
	}

	if !PasswordVerifyWithSalt(password, salt2, hash2) {
		t.Error("hash2 验证失败")
	}

	// 交叉验证应该失败
	if PasswordVerifyWithSalt(password, salt1, hash2) {
		t.Error("使用 salt1 验证 hash2 不应该成功")
	}

	if PasswordVerifyWithSalt(password, salt2, hash1) {
		t.Error("使用 salt2 验证 hash1 不应该成功")
	}
}

// BenchmarkGenerateSalt 基准测试：生成盐值
func BenchmarkGenerateSalt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateSalt()
	}
}

// BenchmarkHashPassword 基准测试：哈希密码
func BenchmarkHashPassword(b *testing.B) {
	password := "test123456"
	salt, _ := GenerateSalt()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashPassword(password, salt)
	}
}

// BenchmarkVerifyPassword 基准测试：验证密码
func BenchmarkVerifyPassword(b *testing.B) {
	password := "test123456"
	salt, _ := GenerateSalt()
	hash, _ := HashPassword(password, salt)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PasswordVerifyWithSalt(password, salt, hash)
	}
}

