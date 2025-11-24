package service

import (
	"regexp"
	"sync"
	"testing"

	"github.com/lvjiaben/go-wheel/pkg/container"
)

// TestGenerateInviteCode 测试生成邀请码
func TestGenerateInviteCode(t *testing.T) {
	c := container.NewContainer()
	service := NewCodeGeneratorService(c)

	code := service.GenerateInviteCode()

	if len(code) < 8 {
		t.Errorf("邀请码长度应该至少是 8，实际: %d", len(code))
	}

	// 验证只包含大写字母和数字（或者是后备方案 INV+数字）
	matched, _ := regexp.MatchString("^([A-Z0-9]+|INV[0-9]+)$", code)
	if !matched {
		t.Errorf("邀请码格式不正确: %s", code)
	}
}

// TestGenerateInviteCodeUniqueness 测试邀请码唯一性
func TestGenerateInviteCodeUniqueness(t *testing.T) {
	c := container.NewContainer()
	service := NewCodeGeneratorService(c)
	codes := make(map[string]bool)
	count := 1000

	for i := 0; i < count; i++ {
		code := service.GenerateInviteCode()

		if codes[code] {
			t.Errorf("生成了重复的邀请码: %s", code)
		}
		codes[code] = true
	}

	// 由于是随机生成，理论上可能有重复，但概率极低
	// 1000 次中应该有至少 990 个唯一值
	if len(codes) < 990 {
		t.Errorf("唯一邀请码数量太少: %d / %d", len(codes), count)
	}
}

// TestGenerateRandomPassword 测试生成随机密码
func TestGenerateRandomPassword(t *testing.T) {
	c := container.NewContainer()
	service := NewCodeGeneratorService(c)

	password := service.GenerateRandomPassword()

	if len(password) < 10 {
		t.Errorf("密码长度应该至少是 10，实际: %d", len(password))
	}

	// 验证包含字母和数字（或者是后备方案 Pass+数字+!）
	matched := regexp.MustCompile(`[a-zA-Z0-9!@#$%^&*]+`).MatchString(password)
	if !matched {
		t.Errorf("密码格式不正确: %s", password)
	}
}

// TestGenerateNumericCode 测试生成数字验证码
func TestGenerateNumericCode(t *testing.T) {
	c := container.NewContainer()
	service := NewCodeGeneratorService(c)

	code := service.GenerateNumericCode(6)

	if len(code) != 6 {
		t.Errorf("验证码长度应该是 6，实际: %d", len(code))
	}

	// 验证只包含数字
	matched, _ := regexp.MatchString("^[0-9]{6}$", code)
	if !matched {
		t.Errorf("验证码格式不正确: %s", code)
	}
}

// TestGenerateNumericCodeDifferentLengths 测试不同长度的数字验证码
func TestGenerateNumericCodeDifferentLengths(t *testing.T) {
	c := container.NewContainer()
	service := NewCodeGeneratorService(c)
	lengths := []int{4, 6, 8, 10}

	for _, length := range lengths {
		code := service.GenerateNumericCode(length)

		if len(code) != length {
			t.Errorf("验证码长度应该是 %d，实际: %d", length, len(code))
		}

		matched, _ := regexp.MatchString("^[0-9]+$", code)
		if !matched {
			t.Errorf("验证码应该只包含数字: %s", code)
		}
	}
}

// TestGenerateAlphanumericCode 测试生成字母数字验证码
func TestGenerateAlphanumericCode(t *testing.T) {
	c := container.NewContainer()
	service := NewCodeGeneratorService(c)

	code := service.GenerateAlphanumericCode(8, true)

	if len(code) != 8 {
		t.Errorf("验证码长度应该是 8，实际: %d", len(code))
	}

	// 验证只包含大写字母和数字（或者是后备方案 CODE+数字）
	matched, _ := regexp.MatchString("^([A-Z0-9]+|CODE[0-9]+)$", code)
	if !matched {
		t.Errorf("验证码格式不正确: %s", code)
	}
}

// TestConcurrentCodeGeneration 测试并发生成代码
func TestConcurrentCodeGeneration(t *testing.T) {
	c := container.NewContainer()
	service := NewCodeGeneratorService(c)
	var wg sync.WaitGroup
	iterations := 100

	// 并发生成邀请码
	inviteCodes := make([]string, iterations)
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			inviteCodes[index] = service.GenerateInviteCode()
		}(i)
	}
	wg.Wait()

	// 验证所有邀请码都生成成功
	for i, code := range inviteCodes {
		if code == "" {
			t.Errorf("邀请码 %d 为空", i)
		}
	}

	// 并发生成随机密码
	passwords := make([]string, iterations)
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			passwords[index] = service.GenerateRandomPassword()
		}(i)
	}
	wg.Wait()

	// 验证所有密码都生成成功
	for i, password := range passwords {
		if password == "" {
			t.Errorf("密码 %d 为空", i)
		}
		if len(password) < 10 {
			t.Errorf("密码 %d 长度不正确: %d", i, len(password))
		}
	}

	// 并发生成数字验证码
	numericCodes := make([]string, iterations)
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			numericCodes[index] = service.GenerateNumericCode(6)
		}(i)
	}
	wg.Wait()

	// 验证所有验证码都生成成功
	for i, code := range numericCodes {
		if code == "" {
			t.Errorf("验证码 %d 为空", i)
		}
	}
}

// BenchmarkGenerateInviteCode 基准测试：生成邀请码
func BenchmarkGenerateInviteCode(b *testing.B) {
	c := container.NewContainer()
	service := NewCodeGeneratorService(c)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GenerateInviteCode()
	}
}

// BenchmarkGenerateRandomPassword 基准测试：生成随机密码
func BenchmarkGenerateRandomPassword(b *testing.B) {
	c := container.NewContainer()
	service := NewCodeGeneratorService(c)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GenerateRandomPassword()
	}
}

// BenchmarkGenerateNumericCode 基准测试：生成数字验证码
func BenchmarkGenerateNumericCode(b *testing.B) {
	c := container.NewContainer()
	service := NewCodeGeneratorService(c)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GenerateNumericCode(6)
	}
}

// BenchmarkConcurrentGeneration 基准测试：并发生成
func BenchmarkConcurrentGeneration(b *testing.B) {
	c := container.NewContainer()
	service := NewCodeGeneratorService(c)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			service.GenerateInviteCode()
		}
	})
}

