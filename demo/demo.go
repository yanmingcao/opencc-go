package main

import (
	"fmt"

	"github.com/byvoid/opencc-go"
	"github.com/byvoid/opencc-go/pkg/conversion"
	"github.com/byvoid/opencc-go/pkg/dict"
	"github.com/byvoid/opencc-go/pkg/segmentation"
)

func main() {
	fmt.Println("=== OpenCC-Go Demo ===")
	fmt.Println()

	// Demo 1: Basic Dictionary and Conversion
	fmt.Println("1. Basic Dictionary Conversion:")
	basicDemo()

	// Demo 2: Multi-value Dictionary (one-to-many mapping)
	fmt.Println("\n2. Multi-value Dictionary (one-to-many):")
	multiValueDemo()

	// Demo 3: Maximum Forward Matching Segmentation
	fmt.Println("\n3. Maximum Forward Matching Segmentation:")
	segmentationDemo()

	// Demo 4: Conversion Chain
	fmt.Println("\n4. Conversion Chain (multi-stage):")
	chainDemo()

	// Demo 5: DictGroup (multiple dictionaries)
	fmt.Println("\n5. DictGroup (composite dictionaries):")
	groupDemo()
}

func basicDemo() {
	// Create a simple dictionary
	lexicon := dict.NewLexicon()
	lexicon.Add(dict.NewStrSingleValueDictEntry("简体", "簡體"))
	lexicon.Add(dict.NewStrSingleValueDictEntry("汉字", "漢字"))
	lexicon.Sort()

	d := dict.NewTextDict(lexicon)

	// Create segmentation and converter
	seg := segmentation.NewMaxMatchSegmentation(d)
	conv := conversion.NewConversion(d)
	chain := conversion.NewConversionChain([]*conversion.Conversion{conv})
	converter := opencc.NewConverter("basic", seg, chain)

	input := "简体汉字"
	result := converter.Convert(input)
	fmt.Printf("  Input:  %s\n", input)
	fmt.Printf("  Output: %s\n", result)
}

func multiValueDemo() {
	// Create dictionary with one-to-many mapping
	lexicon := dict.NewLexicon()
	lexicon.Add(dict.NewStrMultiValueDictEntry("发", []string{"髪", "發"}))
	lexicon.Add(dict.NewStrSingleValueDictEntry("头发", "頭髪"))
	lexicon.Sort()

	d := dict.NewTextDict(lexicon)

	// Create converter
	seg := segmentation.NewMaxMatchSegmentation(d)
	conv := conversion.NewConversion(d)
	chain := conversion.NewConversionChain([]*conversion.Conversion{conv})
	converter := opencc.NewConverter("multi", seg, chain)

	input := "头发"
	result := converter.Convert(input)
	fmt.Printf("  Input:  %s\n", input)
	fmt.Printf("  Output: %s (first value of multi-value entry)\n", result)
}

func segmentationDemo() {
	// Create dictionary with phrases
	lexicon := dict.NewLexicon()
	lexicon.Add(dict.NewStrSingleValueDictEntry("简体中文", "簡體中文"))
	lexicon.Add(dict.NewStrSingleValueDictEntry("中文", "中文"))
	lexicon.Sort()

	d := dict.NewTextDict(lexicon)
	seg := segmentation.NewMaxMatchSegmentation(d)

	input := "简体中文转换"
	segments := seg.Segment(input)

	fmt.Printf("  Input: %s\n", input)
	fmt.Printf("  Segments (%d):\n", segments.Length())
	for i := 0; i < segments.Length(); i++ {
		fmt.Printf("    [%d]: %s\n", i, segments.At(i))
	}
}

func chainDemo() {
	// Create two dictionaries for two-stage conversion
	lexicon1 := dict.NewLexicon()
	lexicon1.Add(dict.NewStrSingleValueDictEntry("a", "b"))
	lexicon1.Sort()
	d1 := dict.NewTextDict(lexicon1)

	lexicon2 := dict.NewLexicon()
	lexicon2.Add(dict.NewStrSingleValueDictEntry("b", "c"))
	lexicon2.Sort()
	d2 := dict.NewTextDict(lexicon2)

	// Create conversion chain: a -> b -> c
	conv1 := conversion.NewConversion(d1)
	conv2 := conversion.NewConversion(d2)
	chain := conversion.NewConversionChain([]*conversion.Conversion{conv1, conv2})

	// Use character segmentation for this demo
	seg := segmentation.NewCharactersSegmentation()
	converter := opencc.NewConverter("chain", seg, chain)

	input := "a"
	result := converter.Convert(input)
	fmt.Printf("  Stage 1: a -> b\n")
	fmt.Printf("  Stage 2: b -> c\n")
	fmt.Printf("  Input:  %s\n", input)
	fmt.Printf("  Output: %s\n", result)
}

func groupDemo() {
	// Create first dictionary
	lexicon1 := dict.NewLexicon()
	lexicon1.Add(dict.NewStrSingleValueDictEntry("简体", "簡體"))
	lexicon1.Sort()
	d1 := dict.NewTextDict(lexicon1)

	// Create second dictionary with additional entries
	lexicon2 := dict.NewLexicon()
	lexicon2.Add(dict.NewStrSingleValueDictEntry("汉字", "漢字"))
	lexicon2.Sort()
	d2 := dict.NewTextDict(lexicon2)

	// Create DictGroup combining both
	group := dict.NewDictGroup([]dict.Dict{d1, d2})

	// Create converter using the group
	seg := segmentation.NewMaxMatchSegmentation(group)
	conv := conversion.NewConversion(group)
	chain := conversion.NewConversionChain([]*conversion.Conversion{conv})
	converter := opencc.NewConverter("group", seg, chain)

	input := "简体汉字"
	result := converter.Convert(input)
	fmt.Printf("  Dictionary 1: 简体 -> 簡體\n")
	fmt.Printf("  Dictionary 2: 汉字 -> 漢字\n")
	fmt.Printf("  Input:  %s\n", input)
	fmt.Printf("  Output: %s\n", result)
}
