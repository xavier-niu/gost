/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gxchan

import (
	"sync"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestUnboundedChan(t *testing.T) {
	ch := NewUnboundedChan(300)

	var count int

	for i := 1; i < 200; i++ {
		ch.In() <- i
	}

	for i := 1; i < 60; i++ {
		v, _ := <-ch.Out()
		count += v.(int)
	}

	assert.Equal(t, 100, ch.queue.Cap())

	for i := 200; i <= 1200; i++ {
		ch.In() <- i
	}
	assert.Equal(t, 1600, ch.queue.Cap())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		var icount int
		for v := range ch.Out() {
			count += v.(int)
			icount++
			if icount == 900 {
				break
			}
		}
	}()

	wg.Wait()

	close(ch.In())

	// buffer should be empty
	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range ch.Out() {
			count += v.(int)
		}
	}()

	wg.Wait()

	assert.Equal(t, 720600, count)

}

func TestUnboundedChanWithQuota(t *testing.T) {
	ch := NewUnboundedChanWithQuota(10, 15)
	assert.Equal(t, 2, cap(ch.in))
	assert.Equal(t, 3, cap(ch.out))
	assert.Equal(t, 4, ch.queue.Cap())
	assert.Equal(t, 0, ch.Len())
	assert.Equal(t, 10, ch.Cap())

	var count int

	for i:=0; i<10; i++ {
		ch.In() <- i
	}

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 14, ch.Cap())
	assert.Equal(t, 10, ch.Len())

	for i:=0; i<10; i++ {
		v, ok := <- ch.Out()
		assert.True(t, ok)
		count += v.(int)
	}

	assert.Equal(t, 45, count)

	for i:=0; i<15; i++ {
		ch.In() <- i
	}

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 15, ch.Cap())
	assert.Equal(t, 15, ch.Len())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		ch.In() <- 15
	}()

	time.Sleep(10 * time.Millisecond)

	for i:=0; i<16; i++ {
		v, ok := <- ch.Out()
		assert.True(t, ok)
		count += v.(int)
	}

	assert.Equal(t, 165, count)

	wg.Wait()
}
