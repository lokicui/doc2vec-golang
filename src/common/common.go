package common

func Max(first int, args ... int) int {
    for _, v := range args {
        if first < v {
            first = v
        }
    }
    return first
}

func Min(first int, args ... int) int {
    for _, v := range args {
        if first > v {
            first = v
        }
    }
    return first
}
