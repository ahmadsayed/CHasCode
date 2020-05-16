package ceas

import "testing"

func setupTestCase(t *testing.T) func(t *testing.T) {
	t.Log("setup test case")
	CreateCluster()
	return func(t *testing.T) {
		DeleteCluster()
	}
}

func setupSubTest(t *testing.T) func(t *testing.T) {
	t.Log("setup sub test")
	return func(t *testing.T) {
		t.Log("teardown sub test")
	}
}

func TestAddtionFirst(t *testing.T) {
	teardownTestCase := setupTestCase(t)
	defer teardownTestCase(t)
	result := 5
	if result != 5 {
		t.Fatalf("expected %v", 5)
	}
}

func TestAddtionSecond(t *testing.T) {

}
