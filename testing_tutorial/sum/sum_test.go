package sum

import "testing"
func TestItns(t *testing.T){
    tt := []struct{
    name string
    sl []int
    res int
    }{
        {"one", []int{1,2,3,4,5}, 15},
        {"nil", nil, 0},
        {"always", []int{4, -4}, 0},
    }
    for _, tc := range tt{
        t.Run(tc.name, func (t *testing.T){
    s := Ints(tc.sl... )
    if s != tc.res{
    t.Fatalf("%s expecting %d found %v" , tc.name, tc.res, s)
    }
})
}
}

func TestDouble(t *testing.T){
    tt := []struct{
    name string
    input , res int
    }{

        {"two", 2, 4},
        {"three", 3, 6},
        {"eighty", 80, 160},
        {"thousand", 1000, 2000},
    }

    for _, tc := range tt {
        t.Run(tc.name, func(t *testing.T){
    if tc.input * 2 != tc.res{
    t.Fatalf("Double of %v shouldn't be %v.",tc.input, tc.res)
    }
})
    }
}
