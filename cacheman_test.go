package cacheman

import "testing"

func TestAll(t *testing.T) {

	// Initialize
	cm := New(0)

	cm.Set("test1", []byte("Contrary to popular belief, Lorem Ipsum is not simply random text."))
	cm.Set("test2", []byte("It has roots in a piece of classical Latin literature from 45 BC,"))
	cm.Set("test3", []byte("making it over 2000 years old. Richard McClintock, a Latin"))
	cm.Set("test4", []byte("professor at Hampden-Sydney College in Virginia, looked up one of"))
	cm.Set("test5", []byte("the more obscure Latin words, consectetur, from a Lorem Ipsum"))
	cm.Set("madam1", []byte("passage, and going through the cites of the word in classical"))
	cm.Set("madam2", []byte("literature, discovered the undoubtable source. Lorem Ipsum"))
	cm.Set("madam3", []byte("comes from sections 1.10.32 and 1.10.33 of \"de Finibus Bonorum"))
	cm.Set("madam4", []byte("et Malorum\" (The Extremes of Good and Evil) by Cicero, written in"))
	cm.Set("madam5", []byte("45 BC. This book is a treatise on the theory of ethics, very popular"))
	cm.Set("unique", []byte("during the Renaissance. The first line of Lorem Ipsum, \"Lorem"))
	cm.Set("duplicate", []byte("ipsum dolor sit amet..\", comes from a line in section 1.10.32."))

	keys := cm.ListKeys()

	// list keys
	for _, k := range keys {
		t.Log(k)
	}

	t.Log("----------------------------------all keys-------------------------------")

	// delete keys with wildcard pattern
	cm.Del("test*")

	// review keys
	keys = cm.ListKeys()
	for _, k := range keys {
		t.Log(k)
	}

	t.Log("-----------------------------all keys after pattern deletion-------------------------------")

	// delete unique key
	cm.Del("unique")
	keys = cm.ListKeys()
	for _, k := range keys {
		t.Log(k)
	}

	t.Log("-----------------------------all keys after unique deletion-------------------------------")

	keys = cm.ListKeys()
	for _, k := range keys {
		v := string(cm.Get([]byte{}, k))
		t.Log(k, ": ", v)
	}
	t.Log("-----------------------------all keys and values remaning-------------------------------")
}
