package colormask

type ColorMask struct {
	Color map[uint]map[uint]uint32
}

func NewColorMask() *ColorMask {
	cm := ColorMask{}
	c := map[uint]map[uint]uint32{}
	cm.Color = c
	return &cm
}

func (c *ColorMask) AddBoxMask(startX, endX, startY, endY uint, color uint32) {
	for i := startX; i <= endX; i++ {
		for j := startY; j < endY; j++ {
			_, ok := c.Color[i]
			if !ok {
				c.Color[i] = map[uint]uint32{}
			}
			c.Color[i][j] = color
		}
	}
}
